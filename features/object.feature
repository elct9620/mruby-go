Feature: Object
  Scenario: I can new a object
    When I execute ruby code:
      """
      Object.new
      """
    Then there should return object
