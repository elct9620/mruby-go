Feature: String
  Scenario: Load string from pool
    When I execute ruby code:
      """
      "Hello World"
      """
    Then there should return string "Hello World"
