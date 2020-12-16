package util

//Common Key constants
const (
	KeyName         = "name"
	KeyAmount       = "amount"
	KeyCurrency     = "currency"
	KeyNotes        = "notes"
	KeyDescription  = "description"
	KeyUnit         = "unit"
	KeyHsnCode      = "hsn_code"
	KeySacCode      = "sac_code"
	KeyTaxRate      = "tax_rate"
	KeyCess         = "cess"
	KeyTaxInclusive = "tax_inclusive"
	KeyDiscount     = "discount"
	KeyFcmToken     = "fcm_token"
)

//Create Order API Key Constants
const (
	KeyMerchantOrderID    = "merchant_order_id"
	KeyPaymentAutoCapture = "payment_auto_capture"
)

//Create VirtualAccount API Key Constants
const (
	KeyOrderID                    = "order_id"
	KeyCollectionMethods          = "collection_methods"
	KeyNotificationMethods        = "notification_method"
	KeyPollForPaymentStatusUpdate = "poll_for_payment_status_update"
)

//Billing Object IDs
const (
	KeyProductID      = "product_id"
	KeyPlanID         = "plan_id"
	KeyCustomerID     = "customer_id"
	KeySubscriptionID = "subscription_id"
	KeyInvoiceID      = "invoice_id"
	KeyInvoiceItemID  = "id"
)

//Create Product API Key Constants
const (
	KeyType      = "type"
	KeyUnitLabel = "unit_label"
)

//Create Customer API Key Constants
const (
	KeyEmail           = "email"
	KeyContactNo       = "contact_no"
	KeyGstin           = "gstin"
	KeyShippingAddress = "shipping_address"
	KeyBillingAddress  = "billing_address"
)

//Create Plan API Key Constants
const (
	KeyFrequency = "frequency"
	KeyInterval  = "interval"
)

//Create Subscription API Key Constants
const (
	KeyBillingCycleCount      = "billing_cycle_count"
	KeyCustomerNotificationBy = "customer_notification_by"
	KeyTrialDuration          = "trial_duration"
	KeyBillingMethod          = "billing_method"
	KeyDueByDays              = "due_by_days"
	KeyUpfrontItems           = "upfront_items"

	KeyBillingCycleEnd = "at_billing_cycle_end"
	KeyStatus          = "status"
)

//Create Invoice Subscription InvoiceItem
const (
	KeyDueDate                   = "due_date"
	KeyDueDateFrom               = "due_date_from"
	KeyDueDateTo                 = "due_date_to"
	KeyNotifyBy                  = "notify_by"
	KeyInvoiceNo                 = "invoice_no"
	KeyLineItems                 = "line_items"
	KeyCustomerNotes             = "customer_notes"
	KeyTermsConditions           = "terms_conditions"
	KeyPlaceOfSupply             = "place_of_supply"
	KeyReceiptNo                 = "receipt_no"
	KeyPartialPaymentMode        = "partial_payment_mode"
	KeyInvoiceCategory           = "invoice_category"
	KeyDestinationBankAccounts   = "destination_bank_accounts"
	KeyGatewayID                 = "gateway_id"
	KeyAccountNumbers            = "account_numbers"
	KeyInvoiceIDs                = "invoice_ids"
	KeyAllowedSourceBankAccounts = "allowed_source_bank_accounts"
	KeyMerchantInvoiceID         = "merchant_invoice_id"
	KeyMerchantInvoiceItemID     = "merchant_invoice_item_id"
	KeyInvoiceStatus             = "status"
)

// CreateTransfer keys
const (
	KeySourceID      = "source_id"
	KeyBeneficiaryID = "beneficiary_id"
	KeyTransfers     = "transfers"
	KeyTransferMode  = "transfer_mode"
)

//Create Invoice Item Key Constants
const (
	KeyQuantity = "quantity"
)

//Create BeneficiaryAccounts Key Constants
const (
	KeyBusinessName                = "business_name"
	KeyBusinessEntityType          = "business_entity_type"
	KeyBankAccountVerificationMode = "bank_account_verification_mode"
	KeyBeneficiaryName             = "beneficiary_name"
	KeyIFSC                        = "ifsc"
	KeyBankAccountNumber           = "bank_account_number"
	KeyAccountType                 = "account_type"
	KeyBankName                    = "bank_name"
)

//API pagination key constants
const (
	KeyCount = "count"
	KeyFrom  = "from"
	KeyTo    = "to"
	KeySkip  = "skip"
)

//Orders API filter/search key constants
const (
	KeyAuthorized = "authorized"
)

//Preferences API key constants
const (
	KeyAccessID = "access_id"
)

// Boolean value constants
const (
	ValueTrue  = "true"
	ValueFalse = "false"
)

//Payouts API key constants
const (
	KeyRemittanceAccountNo            = "remittance_account_no"
	KeyBeneficiaryBankIfsc            = "beneficiary_ifsc"
	KeyBeneficiaryBankAccountNo       = "beneficiary_account_no"
	KeyBeneficiaryBankBeneficiaryName = "beneficiary_name"
	KeyMerchantReferenceID            = "merchant_reference_id"
	KeyPayoutInstrument               = "instrument"
	KeyPayoutMethod                   = "method"
	KeyPayoutStatus                   = "status"
	KeyPurpose                        = "purpose"
	KeyNarration                      = "narration"
)

// Error message constants
const (
	//UnsupportedParamMsg is given as error message when any unsupported parameter is provided
	UnsupportedParamMsg = "One of the request parameters specified in the URL is not supported"
	//InvalidParameterMsg is given as error message when any validation fails for a field in the API side
	InvalidParameterMsg = "The request has invalid parameters"
	//InvalidPostParameterMsg is given as error message when any validation fails for a field in the API side for POST/PUT request
	InvalidPostParameterMsg = "Invalid value provided in field"
	//Missing mandatory field
	MissingMandatoryField = "Required field does not have a value"
	//InvalidHsnSacCodeMsg is given when both HSN and SAC code is passed from API side
	InvalidHsnSacCodeMsg = "This line item cannot be added with both hsn_code and sac_code"
)

const (
	KeyUserId           = "user_id"
	KeyUserUuid         = "user_uuid"
	KeyIdentityProvider = "identity_provider"
	ClientId            = "client_id"
)

const (
	OverrideAutoCapture          = "override_auto_capture"
	PricingPlanTypeID            = "pricing_plan_typeID"
	SettlementIntervalPlanTypeID = "settlement_interval_plan_typeID"
	IsRegisteredForWebhook       = "is_registered_for_webhook"
	IsRegisteredForTransfer      = "is_registered_for_transfer"
	AccountNo                    = "account_no"
	AccountType                  = "account_type"
	IFSC                         = "IFSC"
	BankName                     = "bank_name"
	BankAddress                  = "bank_address"
	BeneficiaryName              = "beneficiary_name"
	Address                      = "address"
	City                         = "city"
	State                        = "state"
	Pin                          = "pin"
)

const (
	KeyMobileProfileID     = "mobile_profile_id"
	KeyNewToken            = "new_fcm_token"
	KeyModifiedAtFrom      = "modified_at_from"
	KeyModifiedAtTo        = "modified_at_to"
	KeyBusinessDisplayName = "business_display_name"
)

const (
	KeyDeviceID           = "device_id"
	KeyIncreaseSequenceID = "increase_sequence_id"
	KeyActiveState        = "active_state"
)

//AssignTeamRoles Key Constants
const (
	KeyRoleID   = "role_id"
	KeyNickname = "nickname"
)

const (
	KeyProfileType = "profile_type"
)

const (
	KeyContentType = "Content-Type"
)

const (
	KeyVendorID = "vendor_id"
)

const (
	KeyRequest          = "request"
	KeyApprove          = "approve"
	KeyVendors          = "vendors"
	KeySourceBanks      = "source_banks"
	KeyDestinationBanks = "destination_banks"
	KeySource           = "source"
	KeyDestination      = "destination"
)

const (
	KeyHasPortalAccess    = "has_portal_access"
	KeyMerchantCustomerID = "merchant_customer_id"
	KeyLabel              = "label"
)

const (
	KeyMerchantList = "merchant_list"
	KeyAmountFrom   = "amount_from"
	KeyAmountTo     = "amount_to"
)

const (
	KeyFilePath = "file_path"
)

const (
	SyncWithSAP          = "SAP"
	KeySapAmountDue      = "amount_due"
	KeySapCompanyCode    = "company_code"
	KeySapCustomerName   = "customer_name"
	KeySapCustomerNumber = "customer_number"
	KeySapDescription    = "description"
	KeySapItem           = "item"
)

const (
	CurrencyINR = "INR"
)

//for payment update notification to SAP
const (
	KeyRecords        = "Records"
	KeyCustomerNumber = "customer_number"
	KeyCustomerName   = "customer_name"
	KeyCompanyCode    = "company_code"
	KeyItem           = "item"
	KeyAmountDue      = "amount_due"
	KeyPaymentAmount  = "payment_amount"
	KeyBankAccount    = "bank_account"
	KeyTransactionRef = "transaction_ref"
)
