// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	providerconfig "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/providerconfig"
	legacyserviceuser "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/user/legacyserviceuser"
	programmaticaccesstoken "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/user/programmaticaccesstoken"
	serviceuser "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/user/serviceuser"
	sessionpolicyattachment "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/user/sessionpolicyattachment"
	user "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/user/user"
	authenticationpolicyattachment "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/userpreview/authenticationpolicyattachment"
	passwordpolicyattachment "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/userpreview/passwordpolicyattachment"
	publickeys "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/userpreview/publickeys"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		providerconfig.Setup,
		legacyserviceuser.Setup,
		programmaticaccesstoken.Setup,
		serviceuser.Setup,
		sessionpolicyattachment.Setup,
		user.Setup,
		authenticationpolicyattachment.Setup,
		passwordpolicyattachment.Setup,
		publickeys.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		providerconfig.SetupGated,
		legacyserviceuser.SetupGated,
		programmaticaccesstoken.SetupGated,
		serviceuser.SetupGated,
		sessionpolicyattachment.SetupGated,
		user.SetupGated,
		authenticationpolicyattachment.SetupGated,
		passwordpolicyattachment.SetupGated,
		publickeys.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
