package constants

import "fmt"

func MockOfEventType(eventType string) string {
	switch eventType {
	case EventTypePerson:
		return MockPerson
	case EventTypeAddress:
		return MockAddress
	case EventTypePhone:
		return MockPhone
	case EventTypeAnalysis:
		return MockCreditAnalysis
	case EventTypeAccount:
		return MockAccount
	case EventTypeCard:
		return MockCard
	case EventTypeLimit:
		return MockLimit
	case EventTypePurchase:
		return MockPurchase
	case EventTypePurchaseApproved:
		return MockPurchaseApproved
	case EventTypePurchaseCanceled:
		return MockPurchaseCanceled
	case EventTypePurchaseCleared:
		return MockPurchaseCleared
	case EventTypePurchaseProcessed:
		return MockPurchaseProcessed
	case EventTypeDeclined:
		return MockDeclined
	case EventTypeStandIn:
		return MockStandIn
	case EventTypeTransaction:
		return MockTransaction
	case EventTypePayment:
		return MockPayment
	case EventTypeAdjustment:
		return MockAdjustment
	case EventTypeStatement:
		return MockStatement
	case EventTypeBankSlip:
		return MockBankSlip
	case EventTypeToken:
		return MockToken
	default:
		return fmt.Sprintf(`{"error": "invalid event type", "givenEvent": "%s"}`, eventType)
	}
}

const (
	MockPerson = `{
  "person_id": 10,
  "person_type": "PF",
  "name": "João da Silva",
  "tax_identification_number": "23471454012",
  "birth_date": "1990-06-27T03:00:00.000Z",
  "gender": null,
  "document_id": "5884532",
  "document_issuer": "SSP",
  "mother": "Maria de Lourdes de Souza Mendes",
  "father": "",
  "marital_status": null,
  "birth_place": "",
  "nationality": "Brasileira",
  "salary": 10000.00,
  "spouse_name": "",
  "spouse_salary": 0.00,
  "graduation_degree": "",
  "occupation": "",
  "bank_code": 0,
  "bank_agency": "0",
  "bank_account": "",
  "email": "francessunderland@conductor.com.br",
  "politically_exposed": false,
  "properties": {
    "cmd_seq": 7,
    "dt_capture": "2020-01-07T21:45:15.577Z",
    "dt_publish": "2020-01-07T21:45:40.218Z",
    "dt_transaction": "2020-01-07T21:45:14.570Z",
    "issuer_id": 1,
    "issuer_name": "Issuer",
    "operation": "upd"
  }
}`
	MockAddress = `{
  "person_id": 10,
  "address_id": 15,
  "address_type_id": 1,
  "address_type_description": "Residencial",
  "street": "teste",
  "number": "267",
  "complement": null,
  "district": "Alphaville",
  "city": "Barueri",
  "state": "SP",
  "country": null,
  "zip_code": "01406000",
  "properties": {
    "cmd_seq": 2,
    "dt_capture": "2020-01-07T15:03:21.897Z",
    "dt_publish": "2020-01-07T15:03:47.120Z",
    "dt_transaction": "2020-01-07T15:03:10.933Z",
    "issuer_id": 1,
    "issuer_name": "Issuer",
    "operation": "upd"
  }
}
`
	MockPhone = `{
  "person_id": 10,
  "phone_id": 30580780,
  "type_id": 1,
  "phone_type": "Residencial",
  "area_code": "083",
  "number": "996040739",
  "extension": "99",
  "active": true,
  "properties": {
    "cmd_seq": 1,
    "dt_capture": "2020-01-07T21:45:15.437Z",
    "dt_publish": "2020-01-07T21:45:40.071Z",
    "dt_transaction": "2020-01-07T21:45:14.580Z",
    "issuer_id": 1,
    "issuer_name": "Issuer",
    "operation": "upd"
  }
}`

	MockCreditAnalysis = `{
  "credit_id": 24917450,
  "person_id": 10,
  "product_id": 3,
  "product_description": "bank",
  "create_date": "2023-04-12T01:01:00.000Z",
  "status_id": 1,
  "status_description": "Aprovada",
  "billing_day": 1,
  "entry_mode": "PIER",
  "credit_score": 0,
  "proven_income": null,
  "notes": "11-04-2023 10:01:26 - PIER - Cadastro de Conta - -contas - Cadatro de Conta",
  "properties": {
    "cmd_seq": 1,
    "dt_capture": "2020-01-07T15:21:36.540Z",
    "dt_publish": "2020-01-07T15:22:00.154Z",
    "dt_transaction": "2020-01-07T15:21:35.130Z",
    "issuer_id": 1,
    "issuer_name": "Issuer",
    "operation": "upd"
  }
}`

	MockAccount = `{
  "account_id": 10,
  "person_id": 10,
  "status_id": 0,
  "status_description": "Normal",
  "create_date": "2023-10-11T18:45:00.000Z",
  "due_day": 1,
  "status_date": "2023-10-11T15:45:00",
  "delivery_address": "Residencial",
  "next_due_date": "",
  "next_real_due_date": "",
  "charge_date": null,
  "properties": {
    "cmd_seq": 9,
    "dt_capture": "2020-01-07T15:03:12.480Z",
    "dt_publish": "2020-01-07T15:03:36.402Z",
    "dt_transaction": "2020-01-07T15:03:03.290Z",
    "issuer_id": 1,
    "issuer_name": "Issuer",
    "operation": "upd"
  },
  "product": {
    "product_id": 16,
    "product_description": "Bank XXXXXX",
    "product_type": "PrePago MasterCard",
    "issuer_bank_number": 336,
    "issuer_branch_number": "1",
    "issuer_account_number": "4207-2"
  }
}`

	MockCard = `{
  "card_id": 10,
  "account_id": 10,
  "person_id": 10,
  "pan": "2184******2121",
  "bin": 230744,
  "name_on_card": "DARLONIA MATOS",
  "status_id": 1,
  "status_description": "Normal",
  "status_allow_approve": true,
  "expiration_date": "2031-06-30T03:00:00.000Z",
  "issue_date": "2022-06-26T00:29:22.833Z",
  "status_upd_date": "2022-07-01T22:33:00.000Z",
  "stage_id": 6,
  "stage_description": "Desbloqueado sem Codigo",
  "embossing_date": "2022-06-26T04:02:00.000Z",
  "embossing_file": null,
  "owner": "Titular",
  "is_temporary_card": false,
  "card_sequence_number": 1,
  "card_hash": "4488516325120342671",
  "brand": "MasterCard",
  "incorrect_password_attempts": 0,
	"incorrect_password_attempts": 0,
    "properties": {
      "cmd_seq": 2,
      "dt_capture": "2020-02-20T15:42:16.530Z",
      "dt_publish": "2020-02-20T15:42:17.146Z",
      "dt_transaction": "2020-02-20T12:55:19.573Z",
      "issuer_id": 1,
      "issuer_name": "Issuer",
      "operation": "ins"
    }
  }`

	MockLimit = `{
  "account_id": 10,
  "credit_limit": 15000.00,
  "installment_limit": 15000.00,
  "tranche_limit": 15000.00,
  "withdrawal_limit": 15000.00,
  "international_purchase_limit": 15000.00,
  "international_withdrawal_limit": 15000.00,
  "maximum_limit": 15000.00,
  "available_balance": 9900.00,
  "installment_available_balance": 10000.00,
  "tranche_available_balance": 10000.00,
  "withdrawal_available_balance": 10000.00,
  "international_available_balance": 10000.00,
  "international_withdrawal_available": 10000.00,
  "available_point": 0.00,
  "properties": {
    "cmd_seq": 8,
    "dt_capture": "2020-01-07T21:39:06.393Z",
    "dt_publish": "2020-01-07T21:39:31.276Z",
    "dt_transaction": "2020-01-07T21:39:05.563Z",
    "issuer_id": 1,
    "issuer_name": "Issuer",
    "operation": "upd"
  }
}`

	MockPurchase = `{
  "purchase_id": 2368554327,
  "account_id": 10,
  "card_id": 10,
  "transaction_type_id": "10457",
  "transaction_type_description": "Compra Débito - MasterCard",
  "purchase_date": "2023-10-11T17:29:00.000Z",
  "amount": 200.0000,
  "source_amount": 200.0000,
  "settlement_amount": 200.0000,
  "amount_with_tax": 200.0000,
  "installment": 1,
  "amount_first_installment": 200.0000,
  "amount_next_installment": 200.0000,
  "tax_percent": 0.000000000,
  "iof_amount": 0.0000,
  "tac_amount": 0.0000,
  "total_tax": 0.0000,
  "authorization_code": "287736",
  "merchant": "SQA*Dock Joao Pessoa         BRA",
  "merchant_id": 1,
  "status": "Pendente",
  "status_id": 3,
  "currency_code": 986,
  "mcc": 5712,
  "exchange_rate": 5.4637,
  "properties": {
      "dt_transaction": "2020-05-19T14:16:40.637",
      "dt_capture": "2020-05-19T14:26:45.583",
      "issuer_id": 1,
      "issuer_name": "Emissor",
      "cmd_seq": 1,
      "operation": "ins"
  }
}`

	MockPurchaseApproved = `{
  "purchase_id": 10,
  "account_id": 10,
  "card_id": 10,
  "product_id": 9,
  "transaction_type_id": "270",
  "transaction_type_description": "A vista sem juros - MasterCard",
  "pan": "2184******2121",
  "origin": "MCD",
  "resolution_origin": "",
  "purchase_date": "2023-11-22T19:15:00.000Z",
  "processing_date": null,
  "amount": 287.3500,
  "source_amount": 287.3500,
  "settlement_amount": 58.3800,
  "amount_with_tax": 287.3500,
  "installment": 1,
  "amount_first_installment": 287.3500,
  "amount_next_installment": 287.3500,
  "tax_percent": 0.000000000,
  "iof_amount": 0.0000,
  "tac_amount": 0.0000,
  "total_tax": 0.0000,
  "authorization_code": "206684",
  "merchant": "Booking.com Brasil Ser Sao Paulo     BRA",
  "merchant_id": 1,
  "mcc": 7011,
  "entry_mode": "0812",
  "nsu": 6,
  "authorization_date": "2023-11-22T16:15:00",
  "transaction_identification": 11140259536,
  "status": "Pendente",
  "status_id": 3,
  "currency_code": 986,
  "exchange_rate": 5.1407,
  "mcc_description": "LODGING - HOTELS, MOTELS, AND RESORTS",
  "mcc_group_id": 40,
  "mcc_group_description": "T&E",
  "is_tokenized_transaction": false,
  "transaction_uuid": "52ff0260-6b51-4675-bc84-a44f4e8390a9",
  "merchant_code": "270695000205031",
  "acquiring_code": 23316,
  "terminal": "FWUNTVZ3",
  "properties":{
    "cmd_seq":1,
    "dt_capture":"2021-02-10T23:50:21.270Z",
    "dt_publish":"2021-02-10T23:50:21.541Z",
    "dt_transaction":"2021-02-10T23:50:19.213Z",
    "issuer_id":1,
    "issuer_name":"Sample",
    "operation":"upd"
  }
}`

	MockPurchaseCanceled = `{
  "purchase_id": 10,
  "account_id": 10,
  "card_id": 10,
  "product_id": 4,
  "transaction_type_id": "270",
  "transaction_type_description": "A vista sem juros - MasterCard",
  "pan": "2184******2121",
  "origin": "MCD",
  "resolution_origin": "BXA",
  "purchase_date": "2018-11-09T13:39:00.000Z",
  "processing_date": null,
  "amount": 30.7700,
  "source_amount": 30.7700,
  "settlement_amount": 8.1700,
  "amount_with_tax": 30.7700,
  "installment": 1,
  "amount_first_installment": 30.7700,
  "amount_next_installment": 30.7700,
  "tax_percent": 0.000000000,
  "iof_amount": 0.0000,
  "tac_amount": 0.0000,
  "total_tax": 0.0000,
  "authorization_code": "177374",
  "merchant": "RAPPI                  SAO PAULO     BRA",
  "merchant_id": 1,
  "mcc": 4214,
  "entry_mode": "0810",
  "nsu": 741160,
  "authorization_date": "2018-11-09T10:39:00",
  "transaction_identification": 430912390078899,
  "status": "Baixada por prazo",
  "status_id": 20,
  "currency_code": 986,
  "exchange_rate": null,
  "mcc_description": "MOTOR FREIGHT CARRIERS AND TRUCKING - LOCAL AND LONG DISTANCE, MOVING AND STORAGE COMPANIES",
  "mcc_group_id": 43,
  "mcc_group_description": "Transportation",
  "is_tokenized_transaction": false,
  "transaction_uuid": "",
  "merchant_code": "",
  "acquiring_code": 20124,
  "terminal": "FWUNTVZ3",
  "properties":{
    "cmd_seq":1,
    "dt_capture":"2021-02-11T00:18:45.751Z",
    "dt_publish":"2021-02-11T00:18:46.686Z",
    "dt_transaction":"2021-02-11T00:18:44.567Z",
    "issuer_id":1,
    "issuer_name":"Sample",
    "operation":"upd"
  }
}`

	MockPurchaseCleared = `{
    "purchase_id": 10,
    "account_id": 10,
    "card_id": 10,
    "product_id": 4,
    "transaction_type_id": "270",
    "transaction_type_description": "A vista sem juros - MasterCard",
    "pan": "2184******2121",
    "origin": "MCD",
    "resolution_origin": "BXA",
    "purchase_date": "2018-11-09T13:39:00.000Z",
    "processing_date": null,
    "amount": 30.7700,
    "source_amount": 30.7700,
    "settlement_amount": 8.1700,
    "amount_with_tax": 30.7700,
    "installment": 1,
    "amount_first_installment": 30.7700,
    "amount_next_installment": 30.7700,
    "tax_percent": 0.000000000,
    "iof_amount": 0.0000,
    "tac_amount": 0.0000,
    "total_tax": 0.0000,
    "authorization_code": "177374",
    "merchant": "RAPPI                  SAO PAULO     BRA",
    "merchant_id": 1,
    "mcc": 4214,
    "entry_mode": "0810",
    "nsu": 741160,
    "authorization_date": "2018-11-09T10:39:00",
    "transaction_identification": 110912391874116,
    "status": "Normal",
    "status_id": 0,
    "currency_code": 986,
    "exchange_rate": null,
    "mcc_description": "MOTOR FREIGHT CARRIERS AND TRUCKING - LOCAL AND LONG DISTANCE, MOVING AND STORAGE COMPANIES",
    "mcc_group_id": 43,
    "mcc_group_description": "Transportation",
    "is_tokenized_transaction": false,
    "transaction_uuid": "",
    "merchant_code": "",
    "acquiring_code": 20124,
    "terminal": "FWUNTVZ3",
    "properties":{
      "cmd_seq":1,
      "dt_capture":"2021-02-11T00:26:22.462Z",
      "dt_publish":"2021-02-11T00:26:22.626Z",
      "dt_transaction":"2021-02-11T00:06:39.917Z",
      "issuer_id":1,
      "issuer_name":"Sample",
      "operation":"upd"
   }
}`

	MockPurchaseProcessed = `{
    "purchase_id": 10000,
    "account_id": 10,
    "card_id": 10,
    "product_id": 2,
    "transaction_type_id": "270",
    "transaction_type_description": "A vista sem juros - MasterCard",
    "pan": "2184******2121",
    "origin": "MCD",
    "resolution_origin": "DTO",
    "purchase_date": "2018-11-14T14:58:00.000Z",
    "processing_date": "2018-11-16T03:00:00.000Z",
    "amount": 1.2000,
    "source_amount": 1.2000,
    "settlement_amount": 0.3200,
    "amount_with_tax": 1.2000,
    "installment": 1,
    "amount_first_installment": 1.2000,
    "amount_next_installment": 1.2000,
    "tax_percent": 0.000000000,
    "iof_amount": 0.0000,
    "tac_amount": 0.0000,
    "total_tax": 0.0000,
    "authorization_code": "216880",
    "merchant": "BESSER PA                SAO PAULO    BR",
    "merchant_id": 3,
    "mcc": 8299,
    "entry_mode": "0071",
    "nsu": 895134,
    "authorization_date": "2018-11-14T11:58:00",
    "transaction_identification": 111413582889513,
    "status": "Processada",
    "status_id": 2,
    "currency_code": 986,
    "exchange_rate": null,
    "mcc_description": "SCHOOLS AND EDUCATIONAL SERVICES, NOT ELSEWHERE CLASSIFIED",
    "mcc_group_id": 15,
    "mcc_group_description": "Educational",
    "is_tokenized_transaction": false,
    "transaction_uuid": "",
    "merchant_code": "",
    "acquiring_code": 6207,
    "terminal": "WW335422",
    "properties":{
      "cmd_seq":1,
      "dt_capture":"2021-02-11T00:26:22.462Z",
      "dt_publish":"2021-02-11T00:26:22.626Z",
      "dt_transaction":"2021-02-11T00:06:39.917Z",
      "issuer_id":1,
      "issuer_name":"Sample",
      "operation":"ins"
   }
}`

	MockDeclined = `{
    "message_id": 3017309594,
    "card_id": 10,
    "account_id": 10,
    "pan": "2184******2121",
    "date": "2023-11-06T20:52:53.837Z",
    "mti": "0100",
    "processing_code": "003000",
    "amount": 41.1700,
    "nsu": "53914",
    "transaction_type": "OFF-US",
    "installments": 1,
    "network": "Mastercard",
    "entry_mode": "0051",
    "merchant": "DROGARIA ARAUJO 164    BELO HORIZONT BRA",
    "mcc": "5912",
    "currency_code": 986,
    "country_code": "076",
    "flag_PIN": true,
    "flag_CVV1": false,
    "flag_CVV2": false,
    "flag_ICVV": false,
    "flag_TOKEN": false,
    "response_code": "63",
    "authorizer_response": "FALCON_NEGADA_FRAUDE",
    "reason": "Transação negada pelo Falcon",
    "incorrect_password_attempts": 0,
    "card_status": "Normal",
    "account_status": "Normal",
    "entry_configuration": "05",
    "institution_code": "012088",
    "data_chip_input": "5F2A0209869A032311069C01009F02060000000041179F10120110A04003220000000000000000000000FF8407A00000000410109F1A0200769F2701809F34034203009F360200589F03060000000000009F3704D0FC2FC2950500000480009F3303E0F8E89F260828C32D2249FA55F082023900",
    "mcc_description": null,
    "mcc_group_id": 21,
    "mcc_group_description": "Health Care",
    "transaction_uuid": "dae925b9-d606-4912-b6d3-41a385fb8336",   
    "properties": {
      "cmd_seq": 2,
      "dt_capture": "2020-07-21T14:31:58.009Z",
      "dt_publish": "2020-07-21T14:41:58.009Z",
      "dt_transaction": "2020-07-21T11:22:58.913Z",
      "issuer_id": 1,
      "issuer_name": "Issuer",
      "operation": "ins"
    }
}`

	MockStandIn = `{
  "message_id": 3079372034,
  "card_id": 10,
  "account_id": 10,
  "date": "2023-11-20T22:56:34.510Z",
  "mti": "0120",
  "processing_code": "203000",
  "amount": 2583.8000,
  "nsu": "835686",
  "transaction_type": "OFF-US",
  "installments": 1,
  "network": "Mastercard",
  "entry_mode": "0012",
  "merchant": "First Car Pomona       JHB           ZAF",
  "mcc": "7512",
  "currency_code": 710,
  "country_code": "710",
  "flag_PIN": false,
  "flag_CVV1": false,
  "flag_CVV2": false,
  "flag_ICVV": false,
  "flag_TOKEN": false,
  "authorization_code": "872177",
  "response_code": "00",
  "entry_configuration": "01",
  "institution_code": "018841",
  "data_chip_input": "",
  "reason": 1020000,
  "mcc_description": "AUTOMOBILE RENTAL AGENCY",
  "mcc_group_id": 40,
  "mcc_group_description": "T&E",
  "transaction_uuid": "1c501203-910d-47cf-8829-3fb784143b33",
  "properties": {
    "cmd_seq": 1,
    "dt_capture": "2020-07-21T15:12:35.021Z",
    "dt_publish": "2020-07-21T15:12:45.021Z",
    "dt_transaction": "2020-07-21T15:12:32.473Z",
    "issuer_id": 1,
    "issuer_name": "Issuer",
    "operation": "ins"
  }
}`

	MockTransaction = `{
  "transaction_id": 10,
  "account_id": 10,
  "card_id": 10,
  "event_id": 101,
  "transaction_date": "2018-11-20T03:00:00.000Z",
  "source_date": "2018-11-13T03:00:00.000Z",
  "transaction_type_id": 3444,
  "transaction_type": "Compra",
  "transaction_type_description": "Compra a Vista MasterCard",
  "amount": 1.2000,
  "current_installments": 0,
  "total_installments": 0,
  "isCredit": false,
  "merchant_id": 1,
  "merchant": "BESSER PARTNERS          SAO PAULO    BR",
  "properties": {
    "cmd_seq": 10,
    "dt_capture": "2020-01-07T21:45:15.373Z",
    "dt_publish": "2020-01-07T21:45:40.012Z",
    "dt_transaction": "2020-01-07T21:45:14.583Z",
    "issuer_id": 1,
    "issuer_name": "Issuer",
    "operation": "upd"
  }
}`

	MockPayment = `{
  "payment_id": 10,
  "account_id": 10,
  "date_payment": "2019-01-25T04:00:00.000Z",
  "date_incoming_payment": "2019-01-25T03:00:00.000Z",
  "amount": 87.0100,
  "status_id": 2,
  "status_description": "Processada",
  "entry_mode": "BANCO",
  "bank_slip_number": "10000059",
  "authorization_code": "",
  "merchant_id": 0,
  "merchant": "",
  "transaction_type_id": "99",
  "transaction_type_description": "Pagamento no Banco",
  "properties": {
    "cmd_seq": 1,
    "dt_capture": "2020-01-07T21:39:06.677Z",
    "dt_publish": "2020-01-07T21:39:31.299Z",
    "dt_transaction": "2020-01-07T21:39:05.520Z",
    "issuer_id": 1,
    "issuer_name": "Issuer",
    "operation": "upd"
  }
}`

	MockAdjustment = `{
  "adjustment_id": 10,
  "account_id": 10,
  "type_adjustment_id": 435,
  "type_adjustment_description": "Baixa Saldo em Análise-Compra Contest MasterCard",
  "create_date": "2023-10-11T17:45:00.000Z",
  "source_date": "2023-10-11T17:45:00.000Z",
  "amount": 52.0000,
  "isCredit": true,
  "status_id": 0,
  "status_description": "Normal",
  "user": "SERVICOS-144531-1",
  "external_establishment_description": null,
  "properties": {
    "cmd_seq": 1,
    "dt_capture": "2020-01-07T15:03:20.923Z",
    "dt_publish": "2020-01-07T15:03:47.069Z",
    "dt_transaction": "2020-01-07T15:03:11.130Z",
    "issuer_id": 1,
    "issuer_name": "Issuer",
    "operation": "upd"
  }
}`

	MockStatement = `{
  "statement_id": 10,
  "account_id": 10,
  "product_id": 3,
  "due_date": "2018-12-20T03:00:00.000Z",
  "real_due_date": "2018-12-20T03:00:00.000Z",
  "previous_amount": 0.0000,
  "total_debits": 0.0000,
  "total_payments": 0.0000,
  "other_credits": 0.0000,
  "total_amount": 10.0000,
  "minimum_amount": 0.0000,
  "is_extractable": false,
  "closing_date": "2018-12-04T03:41:00.000Z",
  "dollar_exchange_value": 3.9800,
  "cet": 0.0000,
  "document_number": 16,
  "delivery_mode": "papel",
  "properties": {
    "cmd_seq": 2,
    "dt_capture": "2020-01-07T21:45:15.373Z",
    "dt_publish": "2020-01-07T21:45:40.023Z",
    "dt_transaction": "2020-01-07T21:45:14.543Z",
    "issuer_id": 1,
    "issuer_name": "Issuer",
    "operation": "upd"
  }
}`

	MockBankSlip = `{
  "bank_slip_id": 10,
  "account_id": 10,
  "issue_date": "2019-01-10T03:00:00.000Z",
  "amount": 404.0800,
  "due_date": "2019-01-21T03:00:00.000Z",
  "bank_slip_number": "10000020",
  "flag_paid": false,
  "bank_slip_type": "BOLETO PRIVATE",
  "bank_code": null,
  "cnab_import_date": null,
  "cnab_return": "Enviado no arquivo de remessa para registro",
  "properties": {
    "cmd_seq": 1,
    "dt_capture": "2020-01-07T15:03:18.250Z",
    "dt_publish": "2020-01-07T15:03:42.042Z",
    "dt_transaction": "2020-01-07T15:03:02.973Z",
    "issuer_id": 1,
    "issuer_name": "Issuer",
    "operation": "upd"
  }
}`

	MockToken = `{
  "card_id": 10,
  "token_id": 10,
  "reference_id": "DM4MMC00002557546fe9b8adf5d441918465a433a2b71492",
  "pan_reference_id": "V-5420241031110747475999",
  "requestor_id": "50194719390",
  "token": "5762",
  "status": 2,
  "type": "03",
  "expiration_date": "2023-05-31T04:00:00.000Z",
  "created_at": "2020-04-08T18:05:53.357Z",
  "properties":{
    "cmd_seq":1,
    "dt_capture":"2020-12-28T14:51:25.260Z",
    "dt_publish":"2020-12-28T14:51:25.534Z",
    "dt_transaction":"2020-12-28T14:51:24.173Z",
    "issuer_id":1,
    "issuer_name":"Sample",
    "operation":"upd"
  }
}
`
)
