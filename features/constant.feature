Feature: Constant
  Scenario: I can get non existing constant
    When I execute ruby code:
      """
      MY_CONSTANT
      """
    Then there should return nil
