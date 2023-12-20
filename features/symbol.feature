Feature: Symbol
  Scenario: Load string from pool
    When I execute ruby code:
      """
      :ruby
      """
    Then there should return symbol "ruby"
