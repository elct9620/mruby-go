Feature: Object
  Scenario: I can get Object class
    When I execute ruby code:
      """
      Object
      """
    Then there should return object

  Scenario: I can get Module class
    When I execute ruby code:
      """
      Module
      """
    Then there should return object

  Scenario: I can get Class class
    When I execute ruby code:
      """
      Class
      """
    Then there should return object

  Scenario: I can new a object
    When I execute ruby code:
      """
      Object.new
      """
    Then there should return object
