Feature: Class
  Scenario: I can define a class
    When I execute ruby code:
      """
      class Foo
      end
      Foo
      """
    Then there should return class

