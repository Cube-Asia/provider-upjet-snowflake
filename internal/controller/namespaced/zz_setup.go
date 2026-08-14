// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	accountauthenticationpolicyattachment "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/accountauthenticationpolicyattachment"
	accountpasswordpolicyattachment "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/accountpasswordpolicyattachment"
	alert "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/alert"
	apiintegration "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/apiintegration"
	apiintegrationamazonapigateway "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/apiintegrationamazonapigateway"
	apiintegrationazureapimanagement "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/apiintegrationazureapimanagement"
	apiintegrationexternalmcpdynamicclient "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/apiintegrationexternalmcpdynamicclient"
	apiintegrationexternalmcpoauth2 "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/apiintegrationexternalmcpoauth2"
	apiintegrationgitrepositorygithubapp "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/apiintegrationgitrepositorygithubapp"
	apiintegrationgitrepositoryoauth2 "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/apiintegrationgitrepositoryoauth2"
	apiintegrationgitrepositoryprivatelink "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/apiintegrationgitrepositoryprivatelink"
	apiintegrationgitrepositorytoken "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/apiintegrationgitrepositorytoken"
	apiintegrationgooglecloudapigateway "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/apiintegrationgooglecloudapigateway"
	cortexagent "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/cortexagent"
	cortexsearchservice "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/cortexsearchservice"
	dynamictable "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/dynamictable"
	emailnotificationintegration "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/emailnotificationintegration"
	externalfunction "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/externalfunction"
	externaltable "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/externaltable"
	failovergroup "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/failovergroup"
	fileformat "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/fileformat"
	functionjava "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/functionjava"
	functionjavascript "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/functionjavascript"
	functionpython "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/functionpython"
	functionscala "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/functionscala"
	functionsql "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/functionsql"
	icebergtablefromdeltafiles "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/icebergtablefromdeltafiles"
	icebergtablefromfiles "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/icebergtablefromfiles"
	jobservice "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/jobservice"
	managedaccount "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/managedaccount"
	materializedview "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/materializedview"
	networkpolicyattachment "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/networkpolicyattachment"
	notebook "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/notebook"
	notificationintegration "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/notificationintegration"
	objectparameter "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/objectparameter"
	pipe "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/pipe"
	postgresinstance "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/postgresinstance"
	procedurejava "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/procedurejava"
	procedurejavascript "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/procedurejavascript"
	procedurepython "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/procedurepython"
	procedurescala "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/procedurescala"
	proceduresql "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/proceduresql"
	semanticview "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/semanticview"
	sequence "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/sequence"
	share "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/share"
	stage "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/stage"
	storageintegration "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/storageintegration"
	storagelifecyclepolicy "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/storagelifecyclepolicy"
	table "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/table"
	tablecolumnmaskingpolicyapplication "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/tablecolumnmaskingpolicyapplication"
	tableconstraint "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/tableconstraint"
	tablestoragelifecyclepolicyattachment "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/tablestoragelifecyclepolicyattachment"
	userauthenticationpolicyattachment "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/userauthenticationpolicyattachment"
	userpasswordpolicyattachment "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/userpasswordpolicyattachment"
	userpublickeys "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/userpublickeys"
	warehouseadaptive "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/preview/warehouseadaptive"
	providerconfig "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/providerconfig"
	account "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/account"
	accountparameter "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/accountparameter"
	accountrole "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/accountrole"
	accountsessionpolicyattachment "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/accountsessionpolicyattachment"
	apiauthenticationintegrationwithauthorizationcodegrant "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/apiauthenticationintegrationwithauthorizationcodegrant"
	apiauthenticationintegrationwithclientcredentials "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/apiauthenticationintegrationwithclientcredentials"
	apiauthenticationintegrationwithjwtbearer "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/apiauthenticationintegrationwithjwtbearer"
	authenticationpolicy "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/authenticationpolicy"
	catalogintegrationawsglue "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/catalogintegrationawsglue"
	catalogintegrationicebergrest "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/catalogintegrationicebergrest"
	catalogintegrationobjectstorage "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/catalogintegrationobjectstorage"
	catalogintegrationopencatalog "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/catalogintegrationopencatalog"
	computepool "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/computepool"
	currentaccount "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/currentaccount"
	currentorganizationaccount "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/currentorganizationaccount"
	database "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/database"
	databaserole "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/databaserole"
	execute "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/execute"
	externaloauthintegration "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/externaloauthintegration"
	externalvolume "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/externalvolume"
	gitrepository "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/gitrepository"
	grantaccountrole "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/grantaccountrole"
	grantapplicationrole "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/grantapplicationrole"
	grantdatabaserole "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/grantdatabaserole"
	grantownership "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/grantownership"
	grantprivilegestoaccountrole "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/grantprivilegestoaccountrole"
	grantprivilegestodatabaserole "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/grantprivilegestodatabaserole"
	grantprivilegestoshare "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/grantprivilegestoshare"
	imagerepository "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/imagerepository"
	legacyserviceuser "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/legacyserviceuser"
	listing "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/listing"
	maskingpolicy "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/maskingpolicy"
	networkpolicy "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/networkpolicy"
	networkrule "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/networkrule"
	oauthintegrationforcustomclients "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/oauthintegrationforcustomclients"
	oauthintegrationforpartnerapplications "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/oauthintegrationforpartnerapplications"
	passwordpolicy "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/passwordpolicy"
	primaryconnection "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/primaryconnection"
	resourcemonitor "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/resourcemonitor"
	rowaccesspolicy "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/rowaccesspolicy"
	saml2integration "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/saml2integration"
	schema "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/schema"
	scimintegration "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/scimintegration"
	secondaryconnection "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/secondaryconnection"
	secondarydatabase "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/secondarydatabase"
	secretwithauthorizationcodegrant "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/secretwithauthorizationcodegrant"
	secretwithbasicauthentication "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/secretwithbasicauthentication"
	secretwithclientcredentials "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/secretwithclientcredentials"
	secretwithgenericstring "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/secretwithgenericstring"
	service "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/service"
	serviceuser "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/serviceuser"
	sessionpolicy "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/sessionpolicy"
	shareddatabase "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/shareddatabase"
	stageexternalazure "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/stageexternalazure"
	stageexternalgcs "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/stageexternalgcs"
	stageexternals3 "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/stageexternals3"
	stageexternals3compatible "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/stageexternals3compatible"
	stageinternal "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/stageinternal"
	storageintegrationaws "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/storageintegrationaws"
	storageintegrationazure "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/storageintegrationazure"
	storageintegrationgcs "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/storageintegrationgcs"
	streamlit "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/streamlit"
	streamondirectorytable "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/streamondirectorytable"
	streamonexternaltable "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/streamonexternaltable"
	streamontable "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/streamontable"
	streamonview "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/streamonview"
	tag "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/tag"
	tagassociation "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/tagassociation"
	task "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/task"
	user "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/user"
	userprogrammaticaccesstoken "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/userprogrammaticaccesstoken"
	usersessionpolicyattachment "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/usersessionpolicyattachment"
	view "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/view"
	warehouse "github.com/Cube-Asia/provider-upjet-snowflake/internal/controller/namespaced/stable/warehouse"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		accountauthenticationpolicyattachment.Setup,
		accountpasswordpolicyattachment.Setup,
		alert.Setup,
		apiintegration.Setup,
		apiintegrationamazonapigateway.Setup,
		apiintegrationazureapimanagement.Setup,
		apiintegrationexternalmcpdynamicclient.Setup,
		apiintegrationexternalmcpoauth2.Setup,
		apiintegrationgitrepositorygithubapp.Setup,
		apiintegrationgitrepositoryoauth2.Setup,
		apiintegrationgitrepositoryprivatelink.Setup,
		apiintegrationgitrepositorytoken.Setup,
		apiintegrationgooglecloudapigateway.Setup,
		cortexagent.Setup,
		cortexsearchservice.Setup,
		dynamictable.Setup,
		emailnotificationintegration.Setup,
		externalfunction.Setup,
		externaltable.Setup,
		failovergroup.Setup,
		fileformat.Setup,
		functionjava.Setup,
		functionjavascript.Setup,
		functionpython.Setup,
		functionscala.Setup,
		functionsql.Setup,
		icebergtablefromdeltafiles.Setup,
		icebergtablefromfiles.Setup,
		jobservice.Setup,
		managedaccount.Setup,
		materializedview.Setup,
		networkpolicyattachment.Setup,
		notebook.Setup,
		notificationintegration.Setup,
		objectparameter.Setup,
		pipe.Setup,
		postgresinstance.Setup,
		procedurejava.Setup,
		procedurejavascript.Setup,
		procedurepython.Setup,
		procedurescala.Setup,
		proceduresql.Setup,
		semanticview.Setup,
		sequence.Setup,
		share.Setup,
		stage.Setup,
		storageintegration.Setup,
		storagelifecyclepolicy.Setup,
		table.Setup,
		tablecolumnmaskingpolicyapplication.Setup,
		tableconstraint.Setup,
		tablestoragelifecyclepolicyattachment.Setup,
		userauthenticationpolicyattachment.Setup,
		userpasswordpolicyattachment.Setup,
		userpublickeys.Setup,
		warehouseadaptive.Setup,
		providerconfig.Setup,
		account.Setup,
		accountparameter.Setup,
		accountrole.Setup,
		accountsessionpolicyattachment.Setup,
		apiauthenticationintegrationwithauthorizationcodegrant.Setup,
		apiauthenticationintegrationwithclientcredentials.Setup,
		apiauthenticationintegrationwithjwtbearer.Setup,
		authenticationpolicy.Setup,
		catalogintegrationawsglue.Setup,
		catalogintegrationicebergrest.Setup,
		catalogintegrationobjectstorage.Setup,
		catalogintegrationopencatalog.Setup,
		computepool.Setup,
		currentaccount.Setup,
		currentorganizationaccount.Setup,
		database.Setup,
		databaserole.Setup,
		execute.Setup,
		externaloauthintegration.Setup,
		externalvolume.Setup,
		gitrepository.Setup,
		grantaccountrole.Setup,
		grantapplicationrole.Setup,
		grantdatabaserole.Setup,
		grantownership.Setup,
		grantprivilegestoaccountrole.Setup,
		grantprivilegestodatabaserole.Setup,
		grantprivilegestoshare.Setup,
		imagerepository.Setup,
		legacyserviceuser.Setup,
		listing.Setup,
		maskingpolicy.Setup,
		networkpolicy.Setup,
		networkrule.Setup,
		oauthintegrationforcustomclients.Setup,
		oauthintegrationforpartnerapplications.Setup,
		passwordpolicy.Setup,
		primaryconnection.Setup,
		resourcemonitor.Setup,
		rowaccesspolicy.Setup,
		saml2integration.Setup,
		schema.Setup,
		scimintegration.Setup,
		secondaryconnection.Setup,
		secondarydatabase.Setup,
		secretwithauthorizationcodegrant.Setup,
		secretwithbasicauthentication.Setup,
		secretwithclientcredentials.Setup,
		secretwithgenericstring.Setup,
		service.Setup,
		serviceuser.Setup,
		sessionpolicy.Setup,
		shareddatabase.Setup,
		stageexternalazure.Setup,
		stageexternalgcs.Setup,
		stageexternals3.Setup,
		stageexternals3compatible.Setup,
		stageinternal.Setup,
		storageintegrationaws.Setup,
		storageintegrationazure.Setup,
		storageintegrationgcs.Setup,
		streamlit.Setup,
		streamondirectorytable.Setup,
		streamonexternaltable.Setup,
		streamontable.Setup,
		streamonview.Setup,
		tag.Setup,
		tagassociation.Setup,
		task.Setup,
		user.Setup,
		userprogrammaticaccesstoken.Setup,
		usersessionpolicyattachment.Setup,
		view.Setup,
		warehouse.Setup,
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
		accountauthenticationpolicyattachment.SetupGated,
		accountpasswordpolicyattachment.SetupGated,
		alert.SetupGated,
		apiintegration.SetupGated,
		apiintegrationamazonapigateway.SetupGated,
		apiintegrationazureapimanagement.SetupGated,
		apiintegrationexternalmcpdynamicclient.SetupGated,
		apiintegrationexternalmcpoauth2.SetupGated,
		apiintegrationgitrepositorygithubapp.SetupGated,
		apiintegrationgitrepositoryoauth2.SetupGated,
		apiintegrationgitrepositoryprivatelink.SetupGated,
		apiintegrationgitrepositorytoken.SetupGated,
		apiintegrationgooglecloudapigateway.SetupGated,
		cortexagent.SetupGated,
		cortexsearchservice.SetupGated,
		dynamictable.SetupGated,
		emailnotificationintegration.SetupGated,
		externalfunction.SetupGated,
		externaltable.SetupGated,
		failovergroup.SetupGated,
		fileformat.SetupGated,
		functionjava.SetupGated,
		functionjavascript.SetupGated,
		functionpython.SetupGated,
		functionscala.SetupGated,
		functionsql.SetupGated,
		icebergtablefromdeltafiles.SetupGated,
		icebergtablefromfiles.SetupGated,
		jobservice.SetupGated,
		managedaccount.SetupGated,
		materializedview.SetupGated,
		networkpolicyattachment.SetupGated,
		notebook.SetupGated,
		notificationintegration.SetupGated,
		objectparameter.SetupGated,
		pipe.SetupGated,
		postgresinstance.SetupGated,
		procedurejava.SetupGated,
		procedurejavascript.SetupGated,
		procedurepython.SetupGated,
		procedurescala.SetupGated,
		proceduresql.SetupGated,
		semanticview.SetupGated,
		sequence.SetupGated,
		share.SetupGated,
		stage.SetupGated,
		storageintegration.SetupGated,
		storagelifecyclepolicy.SetupGated,
		table.SetupGated,
		tablecolumnmaskingpolicyapplication.SetupGated,
		tableconstraint.SetupGated,
		tablestoragelifecyclepolicyattachment.SetupGated,
		userauthenticationpolicyattachment.SetupGated,
		userpasswordpolicyattachment.SetupGated,
		userpublickeys.SetupGated,
		warehouseadaptive.SetupGated,
		providerconfig.SetupGated,
		account.SetupGated,
		accountparameter.SetupGated,
		accountrole.SetupGated,
		accountsessionpolicyattachment.SetupGated,
		apiauthenticationintegrationwithauthorizationcodegrant.SetupGated,
		apiauthenticationintegrationwithclientcredentials.SetupGated,
		apiauthenticationintegrationwithjwtbearer.SetupGated,
		authenticationpolicy.SetupGated,
		catalogintegrationawsglue.SetupGated,
		catalogintegrationicebergrest.SetupGated,
		catalogintegrationobjectstorage.SetupGated,
		catalogintegrationopencatalog.SetupGated,
		computepool.SetupGated,
		currentaccount.SetupGated,
		currentorganizationaccount.SetupGated,
		database.SetupGated,
		databaserole.SetupGated,
		execute.SetupGated,
		externaloauthintegration.SetupGated,
		externalvolume.SetupGated,
		gitrepository.SetupGated,
		grantaccountrole.SetupGated,
		grantapplicationrole.SetupGated,
		grantdatabaserole.SetupGated,
		grantownership.SetupGated,
		grantprivilegestoaccountrole.SetupGated,
		grantprivilegestodatabaserole.SetupGated,
		grantprivilegestoshare.SetupGated,
		imagerepository.SetupGated,
		legacyserviceuser.SetupGated,
		listing.SetupGated,
		maskingpolicy.SetupGated,
		networkpolicy.SetupGated,
		networkrule.SetupGated,
		oauthintegrationforcustomclients.SetupGated,
		oauthintegrationforpartnerapplications.SetupGated,
		passwordpolicy.SetupGated,
		primaryconnection.SetupGated,
		resourcemonitor.SetupGated,
		rowaccesspolicy.SetupGated,
		saml2integration.SetupGated,
		schema.SetupGated,
		scimintegration.SetupGated,
		secondaryconnection.SetupGated,
		secondarydatabase.SetupGated,
		secretwithauthorizationcodegrant.SetupGated,
		secretwithbasicauthentication.SetupGated,
		secretwithclientcredentials.SetupGated,
		secretwithgenericstring.SetupGated,
		service.SetupGated,
		serviceuser.SetupGated,
		sessionpolicy.SetupGated,
		shareddatabase.SetupGated,
		stageexternalazure.SetupGated,
		stageexternalgcs.SetupGated,
		stageexternals3.SetupGated,
		stageexternals3compatible.SetupGated,
		stageinternal.SetupGated,
		storageintegrationaws.SetupGated,
		storageintegrationazure.SetupGated,
		storageintegrationgcs.SetupGated,
		streamlit.SetupGated,
		streamondirectorytable.SetupGated,
		streamonexternaltable.SetupGated,
		streamontable.SetupGated,
		streamonview.SetupGated,
		tag.SetupGated,
		tagassociation.SetupGated,
		task.SetupGated,
		user.SetupGated,
		userprogrammaticaccesstoken.SetupGated,
		usersessionpolicyattachment.SetupGated,
		view.SetupGated,
		warehouse.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupWebhookWithManager registers conversion webhooks for all resource kinds in the group.
func SetupWebhookWithManager(mgr ctrl.Manager) error {
	for _, setup := range []func(ctrl.Manager) error{
		accountauthenticationpolicyattachment.SetupWebhookWithManager,
		accountpasswordpolicyattachment.SetupWebhookWithManager,
		alert.SetupWebhookWithManager,
		apiintegration.SetupWebhookWithManager,
		apiintegrationamazonapigateway.SetupWebhookWithManager,
		apiintegrationazureapimanagement.SetupWebhookWithManager,
		apiintegrationexternalmcpdynamicclient.SetupWebhookWithManager,
		apiintegrationexternalmcpoauth2.SetupWebhookWithManager,
		apiintegrationgitrepositorygithubapp.SetupWebhookWithManager,
		apiintegrationgitrepositoryoauth2.SetupWebhookWithManager,
		apiintegrationgitrepositoryprivatelink.SetupWebhookWithManager,
		apiintegrationgitrepositorytoken.SetupWebhookWithManager,
		apiintegrationgooglecloudapigateway.SetupWebhookWithManager,
		cortexagent.SetupWebhookWithManager,
		cortexsearchservice.SetupWebhookWithManager,
		dynamictable.SetupWebhookWithManager,
		emailnotificationintegration.SetupWebhookWithManager,
		externalfunction.SetupWebhookWithManager,
		externaltable.SetupWebhookWithManager,
		failovergroup.SetupWebhookWithManager,
		fileformat.SetupWebhookWithManager,
		functionjava.SetupWebhookWithManager,
		functionjavascript.SetupWebhookWithManager,
		functionpython.SetupWebhookWithManager,
		functionscala.SetupWebhookWithManager,
		functionsql.SetupWebhookWithManager,
		icebergtablefromdeltafiles.SetupWebhookWithManager,
		icebergtablefromfiles.SetupWebhookWithManager,
		jobservice.SetupWebhookWithManager,
		managedaccount.SetupWebhookWithManager,
		materializedview.SetupWebhookWithManager,
		networkpolicyattachment.SetupWebhookWithManager,
		notebook.SetupWebhookWithManager,
		notificationintegration.SetupWebhookWithManager,
		objectparameter.SetupWebhookWithManager,
		pipe.SetupWebhookWithManager,
		postgresinstance.SetupWebhookWithManager,
		procedurejava.SetupWebhookWithManager,
		procedurejavascript.SetupWebhookWithManager,
		procedurepython.SetupWebhookWithManager,
		procedurescala.SetupWebhookWithManager,
		proceduresql.SetupWebhookWithManager,
		semanticview.SetupWebhookWithManager,
		sequence.SetupWebhookWithManager,
		share.SetupWebhookWithManager,
		stage.SetupWebhookWithManager,
		storageintegration.SetupWebhookWithManager,
		storagelifecyclepolicy.SetupWebhookWithManager,
		table.SetupWebhookWithManager,
		tablecolumnmaskingpolicyapplication.SetupWebhookWithManager,
		tableconstraint.SetupWebhookWithManager,
		tablestoragelifecyclepolicyattachment.SetupWebhookWithManager,
		userauthenticationpolicyattachment.SetupWebhookWithManager,
		userpasswordpolicyattachment.SetupWebhookWithManager,
		userpublickeys.SetupWebhookWithManager,
		warehouseadaptive.SetupWebhookWithManager,
		providerconfig.SetupWebhookWithManager,
		account.SetupWebhookWithManager,
		accountparameter.SetupWebhookWithManager,
		accountrole.SetupWebhookWithManager,
		accountsessionpolicyattachment.SetupWebhookWithManager,
		apiauthenticationintegrationwithauthorizationcodegrant.SetupWebhookWithManager,
		apiauthenticationintegrationwithclientcredentials.SetupWebhookWithManager,
		apiauthenticationintegrationwithjwtbearer.SetupWebhookWithManager,
		authenticationpolicy.SetupWebhookWithManager,
		catalogintegrationawsglue.SetupWebhookWithManager,
		catalogintegrationicebergrest.SetupWebhookWithManager,
		catalogintegrationobjectstorage.SetupWebhookWithManager,
		catalogintegrationopencatalog.SetupWebhookWithManager,
		computepool.SetupWebhookWithManager,
		currentaccount.SetupWebhookWithManager,
		currentorganizationaccount.SetupWebhookWithManager,
		database.SetupWebhookWithManager,
		databaserole.SetupWebhookWithManager,
		execute.SetupWebhookWithManager,
		externaloauthintegration.SetupWebhookWithManager,
		externalvolume.SetupWebhookWithManager,
		gitrepository.SetupWebhookWithManager,
		grantaccountrole.SetupWebhookWithManager,
		grantapplicationrole.SetupWebhookWithManager,
		grantdatabaserole.SetupWebhookWithManager,
		grantownership.SetupWebhookWithManager,
		grantprivilegestoaccountrole.SetupWebhookWithManager,
		grantprivilegestodatabaserole.SetupWebhookWithManager,
		grantprivilegestoshare.SetupWebhookWithManager,
		imagerepository.SetupWebhookWithManager,
		legacyserviceuser.SetupWebhookWithManager,
		listing.SetupWebhookWithManager,
		maskingpolicy.SetupWebhookWithManager,
		networkpolicy.SetupWebhookWithManager,
		networkrule.SetupWebhookWithManager,
		oauthintegrationforcustomclients.SetupWebhookWithManager,
		oauthintegrationforpartnerapplications.SetupWebhookWithManager,
		passwordpolicy.SetupWebhookWithManager,
		primaryconnection.SetupWebhookWithManager,
		resourcemonitor.SetupWebhookWithManager,
		rowaccesspolicy.SetupWebhookWithManager,
		saml2integration.SetupWebhookWithManager,
		schema.SetupWebhookWithManager,
		scimintegration.SetupWebhookWithManager,
		secondaryconnection.SetupWebhookWithManager,
		secondarydatabase.SetupWebhookWithManager,
		secretwithauthorizationcodegrant.SetupWebhookWithManager,
		secretwithbasicauthentication.SetupWebhookWithManager,
		secretwithclientcredentials.SetupWebhookWithManager,
		secretwithgenericstring.SetupWebhookWithManager,
		service.SetupWebhookWithManager,
		serviceuser.SetupWebhookWithManager,
		sessionpolicy.SetupWebhookWithManager,
		shareddatabase.SetupWebhookWithManager,
		stageexternalazure.SetupWebhookWithManager,
		stageexternalgcs.SetupWebhookWithManager,
		stageexternals3.SetupWebhookWithManager,
		stageexternals3compatible.SetupWebhookWithManager,
		stageinternal.SetupWebhookWithManager,
		storageintegrationaws.SetupWebhookWithManager,
		storageintegrationazure.SetupWebhookWithManager,
		storageintegrationgcs.SetupWebhookWithManager,
		streamlit.SetupWebhookWithManager,
		streamondirectorytable.SetupWebhookWithManager,
		streamonexternaltable.SetupWebhookWithManager,
		streamontable.SetupWebhookWithManager,
		streamonview.SetupWebhookWithManager,
		tag.SetupWebhookWithManager,
		tagassociation.SetupWebhookWithManager,
		task.SetupWebhookWithManager,
		user.SetupWebhookWithManager,
		userprogrammaticaccesstoken.SetupWebhookWithManager,
		usersessionpolicyattachment.SetupWebhookWithManager,
		view.SetupWebhookWithManager,
		warehouse.SetupWebhookWithManager,
	} {
		if err := setup(mgr); err != nil {
			return err
		}
	}
	return nil
}
