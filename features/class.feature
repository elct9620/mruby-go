Feature: Class
  Scenario: I can define a class
    When I execute ruby code:
      """
      class Foo
      end
      """
    Then there should return class

