Feature: Class
  Scenario: I can get Object class
    When I execute ruby code:
      """
      Object
      """
    Then there should return object
    And there should return class "Object"

  Scenario: I can get Module class
    When I execute ruby code:
      """
      Module
      """
    Then there should return object
    And there should return class "Module"

  Scenario: I can get Class class
    When I execute ruby code:
      """
      Class
      """
    Then there should return object
    And there should return class "Class"

  Scenario: I can get Kernel module
    When I execute ruby code:
      """
      Kernel
      """
    Then there should return object
    And there should return module "Kernel"

  Scenario: I can define a class
    When I execute ruby code:
      """
      class Foo
      end
      Foo
      """
    Then there should return class "Foo"
