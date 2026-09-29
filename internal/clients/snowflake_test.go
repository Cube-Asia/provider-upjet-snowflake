package clients

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Cube-Asia/provider-upjet-snowflake/apis/cluster/v1beta1"
	namespacedv1beta1 "github.com/Cube-Asia/provider-upjet-snowflake/apis/namespaced/v1beta1"

	rfake "github.com/crossplane/crossplane-runtime/v2/pkg/resource/fake"
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
	crfake "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

// fixtureProviderConfig returns a ProviderConfig resolving its credentials
// from a Secret, plus the Secret itself, installed on a fake client.
func fixtureProviderConfig(t *testing.T, credentials []byte) crclient.Client {
	t.Helper()

	scheme := runtime.NewScheme()
	if err := v1beta1.SchemeBuilder.AddToScheme(scheme); err != nil {
		t.Fatalf("cannot add provider scheme: %v", err)
	}
	if err := clientgoscheme.AddToScheme(scheme); err != nil {
		t.Fatalf("cannot add client-go scheme: %v", err)
	}

	pc := &v1beta1.ProviderConfig{
		ObjectMeta: metav1.ObjectMeta{Name: "default"},
		Spec: v1beta1.ProviderConfigSpec{
			Credentials: v1beta1.ProviderCredentials{
				Source: xpv2.CredentialsSourceSecret,
				CommonCredentialSelectors: xpv2.CommonCredentialSelectors{
					SecretRef: &xpv2.SecretKeySelector{
						SecretReference: xpv2.SecretReference{Name: "creds", Namespace: "ns"},
						Key:             "credentials",
					},
				},
			},
		},
	}
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: "creds", Namespace: "ns"},
		Data:       map[string][]byte{"credentials": credentials},
	}

	return crfake.NewClientBuilder().WithScheme(scheme).WithObjects(pc, secret).Build()
}

func fixtureManaged() *rfake.LegacyManaged {
	return &rfake.LegacyManaged{
		ObjectMeta: metav1.ObjectMeta{Name: "xr", UID: "unit-1"},
		LegacyProviderConfigReferencer: rfake.LegacyProviderConfigReferencer{
			Ref: &xpv2.Reference{Name: "default"},
		},
	}
}

func stubConfiguring(rec *atomic.Int32, result any, err error) func(context.Context, *ujconfig.Provider, map[string]any) (any, error) {
	return func(context.Context, *ujconfig.Provider, map[string]any) (any, error) {
		rec.Add(1)
		return result, err
	}
}

// A reconcile stream against an unchanged provider configuration mints the
// Snowflake session exactly once, however many reconciles happen.
func TestTerraformSetupBuilderConfiguresOncePerConfiguration(t *testing.T) {
	kClient := fixtureProviderConfig(t, []byte(`{"account_name":"acct-a","user":"u"}`))
	mg := fixtureManaged()

	probe := &struct{ n int }{n: 1}
	var calls atomic.Int32
	restore := configureAndMeta
	t.Cleanup(func() { configureAndMeta = restore })
	configureAndMeta = stubConfiguring(&calls, probe, nil)

	setup := TerraformSetupBuilder("v1", "source", "v1.2.3", &ujconfig.Provider{})
	ctx := context.Background()

	for i := range 5 {
		ps, err := setup(ctx, kClient, mg)
		if err != nil {
			t.Fatalf("setup call %d: %v", i, err)
		}
		if ps.Meta != probe {
			t.Fatalf("setup call %d: meta should be the memoized probe", i)
		}
		if got := ps.Configuration["account_name"]; got != "acct-a" {
			t.Fatalf("setup call %d: configuration not populated: %v", i, got)
		}
	}
	if got := calls.Load(); got != 1 {
		t.Fatalf("provider configured %d times for an unchanged configuration, want 1", got)
	}
}

// Concurrent reconciles on the same credentials share a single login instead
// of each minting their own session.
func TestTerraformSetupBuilderConcurrentReconcilesMintOneSession(t *testing.T) {
	kClient := fixtureProviderConfig(t, []byte(`{"account_name":"acct-a","user":"u"}`))
	mg := fixtureManaged()

	var calls atomic.Int32
	restore := configureAndMeta
	t.Cleanup(func() { configureAndMeta = restore })
	configureAndMeta = stubConfiguring(&calls, &struct{ n int }{n: 2}, nil)

	setup := TerraformSetupBuilder("v1", "source", "v1.2.3", &ujconfig.Provider{})
	ctx := context.Background()

	const reconciles = 16
	errs := make(chan error, reconciles)
	var wg sync.WaitGroup
	for range reconciles {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ps, err := setup(ctx, kClient, mg)
			if err == nil && ps.Meta == nil {
				err = errors.New("nil meta")
			}
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatalf("concurrent setup: %v", err)
		}
	}
	if got := calls.Load(); got != 1 {
		t.Fatalf("%d concurrent reconciles configured the provider %d times, want 1", reconciles, got)
	}
}

// Changing the provider configuration (e.g. credential rotation through the
// ProviderConfig Secret) invalidates the memoized session.
func TestTerraformSetupBuilderReconfiguresOnConfigurationChange(t *testing.T) {
	kClient := fixtureProviderConfig(t, []byte(`{"account_name":"acct-a","user":"u"}`))
	mg := fixtureManaged()

	var calls atomic.Int32
	restore := configureAndMeta
	t.Cleanup(func() { configureAndMeta = restore })
	configureAndMeta = stubConfiguring(&calls, nil, nil)

	setup := TerraformSetupBuilder("v1", "source", "v1.2.3", &ujconfig.Provider{})
	ctx := context.Background()

	if _, err := setup(ctx, kClient, mg); err != nil {
		t.Fatalf("initial setup: %v", err)
	}

	probe2 := &struct{ n int }{n: 2}
	configureAndMeta = func(_ context.Context, _ *ujconfig.Provider, cfg map[string]any) (any, error) {
		calls.Add(1)
		return probe2, nil
	}
	secret := &corev1.Secret{}
	if err := kClient.Get(ctx, crclient.ObjectKey{Name: "creds", Namespace: "ns"}, secret); err != nil {
		t.Fatalf("cannot read credential secret: %v", err)
	}
	secret.Data["credentials"] = []byte(`{"account_name":"acct-b","user":"u"}`)
	if err := kClient.Update(ctx, secret); err != nil {
		t.Fatalf("cannot update credential secret: %v", err)
	}

	ps, err := setup(ctx, kClient, mg)
	if err != nil {
		t.Fatalf("setup after rotation: %v", err)
	}
	if ps.Meta != probe2 {
		t.Fatal("meta after rotation should come from the fresh configure call")
	}
	if got := calls.Load(); got != 2 {
		t.Fatalf("provider configured %d times after a configuration change, want 2", got)
	}
	if got := ps.Configuration["account_name"]; got != "acct-b" {
		t.Fatalf("stale configuration rendered after rotation: %v", got)
	}
}

// A failing configure is never memoized: the next reconcile retries it.
func TestTerraformSetupBuilderRetriesAfterConfigureError(t *testing.T) {
	kClient := fixtureProviderConfig(t, []byte(`{"account_name":"acct-a","user":"u"}`))
	mg := fixtureManaged()

	var calls atomic.Int32
	restore := configureAndMeta
	t.Cleanup(func() { configureAndMeta = restore })
	configureAndMeta = stubConfiguring(&calls, nil, errors.New("boom"))

	setup := TerraformSetupBuilder("v1", "source", "v1.2.3", &ujconfig.Provider{})
	ctx := context.Background()

	if _, err := setup(ctx, kClient, mg); err == nil {
		t.Fatal("expected the configure error to propagate")
	}
	if _, err := setup(ctx, kClient, mg); err == nil {
		t.Fatal("expected the second reconcile to retry the configure")
	}
	if got := calls.Load(); got != 2 {
		t.Fatalf("configure attempted %d times after an error, want 2", got)
	}
}

// Alternating reconciles across two distinct provider configurations (live
// and shadow credentials) keep independent sessions instead of thrashing a
// single cache slot.
func TestTerraformSetupBuilderAlternatingConfigurationsKeepIndependentSessions(t *testing.T) {
	kClient := fixtureProviderConfig(t, []byte(`{"account_name":"acct-a","user":"u"}`))
	mg := fixtureManaged()

	metaA, metaB := &struct{ n int }{n: 1}, &struct{ n int }{n: 2}
	var calls atomic.Int32
	restore := configureAndMeta
	t.Cleanup(func() { configureAndMeta = restore })
	configureAndMeta = func(_ context.Context, _ *ujconfig.Provider, cfg map[string]any) (any, error) {
		calls.Add(1)
		if cfg["account_name"] == "acct-a" {
			return metaA, nil
		}
		return metaB, nil
	}

	setup := TerraformSetupBuilder("v1", "source", "v1.2.3", &ujconfig.Provider{})
	ctx := context.Background()

	secret := &corev1.Secret{}
	if err := kClient.Get(ctx, crclient.ObjectKey{Name: "creds", Namespace: "ns"}, secret); err != nil {
		t.Fatalf("cannot read credential secret: %v", err)
	}
	update := func(data string) {
		secret.Data["credentials"] = []byte(data)
		if err := kClient.Update(ctx, secret); err != nil {
			t.Fatalf("cannot update credential secret: %v", err)
		}
	}

	update(`{"account_name":"acct-a"}`)
	for range 2 {
		ps, err := setup(ctx, kClient, mg)
		if err != nil || ps.Meta != metaA {
			t.Fatalf("live config: err=%v meta mismatch", err)
		}
		update(`{"account_name":"acct-b"}`)
		ps, err = setup(ctx, kClient, mg)
		if err != nil || ps.Meta != metaB {
			t.Fatalf("shadow config: err=%v meta mismatch", err)
		}
		update(`{"account_name":"acct-a"}`)
	}
	if got := calls.Load(); got != 2 {
		t.Fatalf("alternating reconciles configured the provider %d times, want 2", got)
	}
}

// An entry older than the cache TTL is reconfigured even though the
// configuration digest is unchanged. gosnowflake renews an expired session
// token with the master token, but no path is proven for an expired master
// token itself, so the cache bounds how long any session may live instead of
// relying on renewal. The clock is injected so the test stays fast and
// deterministic.
func TestMetaCacheExpiresEntriesAfterTTL(t *testing.T) {
	now := time.Unix(0, 0)
	c := &metaCache{ttl: time.Hour, grace: 5 * time.Minute, now: func() time.Time { return now }}

	var calls atomic.Int32
	configure := func() (any, error) {
		n := calls.Add(1)
		return &struct{ n int }{n: int(n)}, nil
	}

	// The first call mints and the immediate second call is memoized.
	first, err := c.metaFor("key", time.Hour, configure)
	if err != nil {
		t.Fatalf("first metaFor: %v", err)
	}
	again, err := c.metaFor("key", time.Hour, configure)
	if err != nil {
		t.Fatalf("second metaFor: %v", err)
	}
	if again != first {
		t.Fatal("entry should be memoized within the TTL")
	}

	// Just inside the TTL the same entry is still served.
	now = now.Add(time.Hour - time.Second)
	inside, err := c.metaFor("key", time.Hour, configure)
	if err != nil {
		t.Fatalf("metaFor inside TTL: %v", err)
	}
	if inside != first {
		t.Fatal("entry one second before expiry should still be served")
	}

	// Past the TTL a fresh session is minted.
	now = now.Add(2 * time.Second)
	expired, err := c.metaFor("key", time.Hour, configure)
	if err != nil {
		t.Fatalf("metaFor past TTL: %v", err)
	}
	if expired == first {
		t.Fatal("entry past the TTL should have been reconfigured")
	}

	// And the fresh entry is memoized in turn.
	fresh, err := c.metaFor("key", time.Hour, configure)
	if err != nil {
		t.Fatalf("metaFor after expiry: %v", err)
	}
	if fresh != expired {
		t.Fatal("fresh entry should be memoized")
	}

	if got := calls.Load(); got != 2 {
		t.Fatalf("provider configured %d times across a TTL boundary, want 2", got)
	}
}

// The cache key is insensitive to map insertion order but sensitive to
// content, so equivalent configurations share a session and changed
// credentials do not.
func TestHashConfigurationIsStable(t *testing.T) {
	first := map[string]any{"user": "u", "account_name": "acct-a", "params": map[string]string{"role": "r"}}
	second := map[string]any{"params": map[string]string{"role": "r"}, "account_name": "acct-a", "user": "u"}
	rotated := map[string]any{"user": "u", "account_name": "acct-b", "params": map[string]string{"role": "r"}}

	firstHash, err := hashConfiguration(first)
	if err != nil {
		t.Fatalf("hash first: %v", err)
	}
	secondHash, err := hashConfiguration(second)
	if err != nil {
		t.Fatalf("hash second: %v", err)
	}
	rotatedHash, err := hashConfiguration(rotated)
	if err != nil {
		t.Fatalf("hash rotated: %v", err)
	}

	if firstHash != secondHash {
		t.Fatal("equivalent configurations hashed differently")
	}
	if firstHash == rotatedHash {
		t.Fatal("different configurations hashed identically")
	}
}

// A sweep past ttl+grace closes the evicted session and deletes its entry,
// whether or not the same key was reconfigured: this is what bounds both
// leaks - the replaced-on-expiry session and the entry of a configuration
// that stopped being resolved (for example after a key rotation). Entries
// inside the grace window stay, so a reconcile still holding the previous
// meta keeps a live pool.
func TestMetaCacheSweepClosesAndDeletesStaleEntries(t *testing.T) {
	now := time.Unix(0, 0)
	c := &metaCache{ttl: time.Hour, grace: 5 * time.Minute, now: func() time.Time { return now }}
	configure := func() (any, error) { return &struct{}{}, nil }

	var closed atomic.Int32
	restore := closeMeta
	t.Cleanup(func() { closeMeta = restore })
	closeMeta = func(any) { closed.Add(1) }

	if _, err := c.metaFor("rotate-away", time.Hour, configure); err != nil {
		t.Fatalf("metaFor rotate-away: %v", err)
	}
	if _, err := c.metaFor("key", time.Hour, configure); err != nil {
		t.Fatalf("metaFor key: %v", err)
	}

	// At ttl+1s the expired "key" entry is replaced but NOT yet closed:
	// the grace window keeps the old pool alive past replacement.
	now = now.Add(time.Hour + time.Second)
	if _, err := c.metaFor("key", time.Hour, configure); err != nil {
		t.Fatalf("metaFor past ttl: %v", err)
	}
	if got := closed.Load(); got != 0 {
		t.Fatalf("closed %d sessions inside the grace window, want 0", got)
	}

	// Past the grace period the superseded "key" meta (retired at 3601s)
	// is closed, and the never-resolved-again "rotate-away" entry leaves
	// the map to await its own close.
	now = now.Add(5 * time.Minute)
	if _, err := c.metaFor("key", time.Hour, configure); err != nil {
		t.Fatalf("metaFor past grace: %v", err)
	}
	if got := closed.Load(); got != 1 {
		t.Fatalf("closed %d sessions past the grace window, want 1", got)
	}
	c.mu.Lock()
	_, hasRotate := c.entries["rotate-away"]
	fresh := c.entries["key"]
	c.mu.Unlock()
	if hasRotate {
		t.Fatal("stale rotate-away entry should have been swept from the map")
	}
	if !fresh.created.After(time.Unix(0, 0)) {
		t.Fatal("the live key entry should have survived the sweep")
	}

	// One more grace period later the rotate-away session is closed too.
	now = now.Add(5 * time.Minute)
	if _, err := c.metaFor("key", time.Hour, configure); err != nil {
		t.Fatalf("metaFor after rotate-away grace: %v", err)
	}
	if got := closed.Load(); got != 2 {
		t.Fatalf("closed %d sessions after the final grace window, want 2", got)
	}
}

// sessionCacheTTL resolves spec.sessionCacheTtl when set and positive, and
// falls back to the default otherwise, so a ProviderConfig controls its own
// re-login rate without a controller restart.
func TestSessionCacheTTLResolves(t *testing.T) {
	hour := time.Hour
	half := 30 * time.Minute
	zero := time.Duration(0)
	for name, tc := range map[string]struct {
		spec *namespacedv1beta1.ProviderConfigSpec
		want time.Duration
	}{
		"nil spec":          {nil, defaultSessionCacheTTL},
		"nil field":         {&namespacedv1beta1.ProviderConfigSpec{}, defaultSessionCacheTTL},
		"set":               {&namespacedv1beta1.ProviderConfigSpec{SessionCacheTTL: &metav1.Duration{Duration: half}}, half},
		"zero means unset":  {&namespacedv1beta1.ProviderConfigSpec{SessionCacheTTL: &metav1.Duration{Duration: zero}}, defaultSessionCacheTTL},
		"negative is unset": {&namespacedv1beta1.ProviderConfigSpec{SessionCacheTTL: &metav1.Duration{Duration: -hour}}, defaultSessionCacheTTL},
	} {
		t.Run(name, func(t *testing.T) {
			if got := sessionCacheTTL(tc.spec); got != tc.want {
				t.Fatalf("sessionCacheTTL = %v, want %v", got, tc.want)
			}
		})
	}
}

// closeMeta reaches the Snowflake SDK client by walking the vendored meta
// type's exported Client field by reflection. That type lives in an internal
// package of the vendored provider, so a field rename there would turn the
// sweep's Close into a silent no-op and reintroduce the leak. This guard
// reads the vendored source and fails when the shape closeMeta relies on
// changes. The check never skips: go.mod's replace directive makes the
// package unbuildable without the vendored checkout, so an absent file
// means a broken environment, not an unrelated one.
func TestCloseMetaGuardAgainstFieldRename(t *testing.T) {
	// go test runs with the package directory as cwd; walk up to the
	// module root before joining the vendored path.
	root, err := os.Getwd()
	if err != nil {
		t.Fatalf("cannot get working directory: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(root, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(root)
		if parent == root {
			t.Fatal("cannot find module root")
		}
		root = parent
	}
	rel := filepath.Join(root, ".work", "snowflakedb", "snowflake", "pkg", "internal", "provider", "provider_context.go")
	src, err := os.ReadFile(rel)
	if err != nil {
		t.Fatalf("cannot read vendored provider context at %s (provisioned by make fetch-snowflake-provider-src): %v", rel, err)
	}
	if !regexp.MustCompile(`Client\s+\*sdk\.Client`).Match(src) {
		t.Fatal("vendored provider.Context no longer has an exported `Client *sdk.Client` field; closeMeta's reflection walk is a silent no-op and must be updated")
	}
}

// A failed configure must not queue the superseded meta for closing twice:
// the old entry stays in the map across the failure, and only a successful
// reconfigure retires it. Double-retirement would close one session twice.
func TestMetaCacheConfigureErrorDoesNotDoubleRetire(t *testing.T) {
	now := time.Unix(0, 0)
	c := &metaCache{ttl: time.Hour, grace: 5 * time.Minute, now: func() time.Time { return now }}
	configure := func() (any, error) { return &struct{}{}, nil }

	var closed atomic.Int32
	restore := closeMeta
	t.Cleanup(func() { closeMeta = restore })
	closeMeta = func(any) { closed.Add(1) }

	if _, err := c.metaFor("key", time.Hour, configure); err != nil {
		t.Fatalf("metaFor: %v", err)
	}

	// After expiry, two consecutive failed reconfigures keep the same old
	// entry in the map; neither may retire it.
	now = now.Add(2 * time.Hour)
	boom := errors.New("boom")
	failing := func() (any, error) { return nil, boom }
	if _, err := c.metaFor("key", time.Hour, failing); !errors.Is(err, boom) {
		t.Fatalf("first failing metaFor: %v", err)
	}
	if _, err := c.metaFor("key", time.Hour, failing); !errors.Is(err, boom) {
		t.Fatalf("second failing metaFor: %v", err)
	}

	// The successful reconfigure retires the old meta exactly once.
	if _, err := c.metaFor("key", time.Hour, configure); err != nil {
		t.Fatalf("metaFor after failures: %v", err)
	}
	now = now.Add(6 * time.Minute)
	if _, err := c.metaFor("key", time.Hour, configure); err != nil {
		t.Fatalf("metaFor past grace: %v", err)
	}
	if got := closed.Load(); got != 1 {
		t.Fatalf("closed %d sessions, want exactly 1", got)
	}
}

// Each entry ages by the TTL in force when it was minted: a reconcile
// against a short-TTL ProviderConfig must not evict a long-TTL
// configuration's still-live session.
func TestMetaCachePerEntryTTLIsolation(t *testing.T) {
	now := time.Unix(0, 0)
	c := &metaCache{ttl: time.Hour, grace: 5 * time.Minute, now: func() time.Time { return now }}

	var calls atomic.Int32
	configure := func() (any, error) {
		n := calls.Add(1)
		return &struct{ n int }{n: int(n)}, nil
	}

	long, err := c.metaFor("long", time.Hour, configure)
	if err != nil {
		t.Fatalf("metaFor long: %v", err)
	}

	// Halfway through the long entry's life, a short-TTL config reconciles
	// repeatedly - well past the short TTL plus grace.
	now = now.Add(30 * time.Minute)
	shortTTL := time.Minute
	for i := 0; i < 3; i++ {
		if _, err := c.metaFor("short", shortTTL, configure); err != nil {
			t.Fatalf("metaFor short %d: %v", i, err)
		}
		now = now.Add(90 * time.Second)
	}
	// That is ~3m45s of short-TTL activity; the long entry (30m old, its
	// own ttl is 1h) must still be served from cache.
	got, err := c.metaFor("long", time.Hour, configure)
	if err != nil {
		t.Fatalf("metaFor long again: %v", err)
	}
	if got != long {
		t.Fatal("short-TTL activity must not evict a live long-TTL entry")
	}
	if got := calls.Load(); got != 4 {
		t.Fatalf("configured %d times, want 4 (long + 3 short)", got)
	}
}

// Editing spec.sessionCacheTtl applies to the live session immediately:
// the entry is re-aged with the new TTL on the next hit, so a lowered one
// expires the session on the next boundary instead of at the old deadline.
func TestMetaCacheEditedTTLAppliesToLiveSession(t *testing.T) {
	now := time.Unix(0, 0)
	c := &metaCache{grace: 5 * time.Minute, now: func() time.Time { return now }}

	var calls atomic.Int32
	configure := func() (any, error) {
		n := calls.Add(1)
		return &struct{ n int }{n: int(n)}, nil
	}

	if _, err := c.metaFor("key", time.Hour, configure); err != nil {
		t.Fatalf("metaFor: %v", err)
	}

	// 90 seconds in - far inside the original hour, but past a lowered
	// one-minute TTL - the next reconcile reconfigures.
	now = now.Add(90 * time.Second)
	got, err := c.metaFor("key", time.Minute, configure)
	if err != nil {
		t.Fatalf("metaFor with lowered ttl: %v", err)
	}
	if calls.Load() != 2 {
		t.Fatalf("lowered TTL did not expire the live session at its boundary (calls=%d)", calls.Load())
	}
	if got == nil {
		t.Fatal("expected a fresh meta")
	}
}
