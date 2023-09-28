Feature: Local Variable
  Scenario: Can assign a local variable
    When I execute ruby code:
      """
      name = "World"
      "Hello #{name}"
      """
    Then there should return string "Hello World"
