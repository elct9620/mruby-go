Feature: String
  Scenario: Load string from pool
    When I execute ruby code:
      """
      "Hello World"
      """
    Then there should return string "Hello World"

  Scenario: Concat string with integer
    When I execute ruby code:
      """
      var = 1
      "Hello #{var}"
      """
    Then there should return string "Hello 1"

  Scenario: Concat string with boolean
    When I execute ruby code:
      """
      var = true
      "Hello #{var}"
      """
    Then there should return string "Hello true"
