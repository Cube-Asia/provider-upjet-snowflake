package clients

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Cube-Asia/provider-upjet-snowflake/apis/cluster/v1beta1"

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
	c := &metaCache{ttl: time.Hour, now: func() time.Time { return now }}

	var calls atomic.Int32
	configure := func() (any, error) {
		n := calls.Add(1)
		return &struct{ n int }{n: int(n)}, nil
	}

	// The first call mints and the immediate second call is memoized.
	first, err := c.metaFor("key", configure)
	if err != nil {
		t.Fatalf("first metaFor: %v", err)
	}
	again, err := c.metaFor("key", configure)
	if err != nil {
		t.Fatalf("second metaFor: %v", err)
	}
	if again != first {
		t.Fatal("entry should be memoized within the TTL")
	}

	// Just inside the TTL the same entry is still served.
	now = now.Add(time.Hour - time.Second)
	inside, err := c.metaFor("key", configure)
	if err != nil {
		t.Fatalf("metaFor inside TTL: %v", err)
	}
	if inside != first {
		t.Fatal("entry one second before expiry should still be served")
	}

	// Past the TTL a fresh session is minted.
	now = now.Add(2 * time.Second)
	expired, err := c.metaFor("key", configure)
	if err != nil {
		t.Fatalf("metaFor past TTL: %v", err)
	}
	if expired == first {
		t.Fatal("entry past the TTL should have been reconfigured")
	}

	// And the fresh entry is memoized in turn.
	fresh, err := c.metaFor("key", configure)
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
