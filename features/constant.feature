Feature: Constant
  Scenario: I can get non existing constant
    When I execute ruby code:
      """
      MY_CONSTANT
      """
    Then there should return nil

  Scenario: I can get existing constant
    When I execute ruby code:
      """
      MY_CONSTANT = 1
      MY_CONSTANT
      """
    Then there should return integer 1
