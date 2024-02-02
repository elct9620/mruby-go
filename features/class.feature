Feature: Class
  Scenario: I can get Kernel module
    When I execute ruby code:
      """
      Kernel
      """
    Then there should return module "Kernel"

  Scenario: I can define a class
    When I execute ruby code:
      """
      class Foo
      end
      Foo
      """
    Then there should return class "Foo"
