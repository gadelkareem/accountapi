Feature: Testing Account API
  Services should be able to use Form3 Account API

  Scenario: Register an existing bank account with Form3 or create a new one
    Given Form3 API server is running on "http://accountapi:8080/"
    When services create an account using data
    """
    {
          "id": "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
          "organisation_id": "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
          "country": "GB",
          "base_currency": "GBP",
          "account_number": "41426819",
          "bank_id": "400300",
          "bank_id_code": "GBDSC",
          "bic": "NWBKGB22",
          "iban": "GB11NWBK40030041426819",
          "title": "Ms",
          "first_name": "Samantha",
          "bank_account_name": "Samantha Holder",
          "alternative_bank_account_names": [
            "Sam Holder"
          ],
          "account_classification": "Personal",
          "joint_account": false,
          "account_matching_opt_out": false,
          "secondary_identification": "A1B2C3D4"
    }
    """
    Then the server saves and returns the account in the response

  Scenario: Get a single account using the account ID
    Given Form3 API server is running on "http://accountapi:8080/"
    When services fetch an account using account ID "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"
    Then the server returns the account with ID "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"

