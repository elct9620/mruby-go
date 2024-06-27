Feature: Module
  Scenario: I can define a Module
    When I execute ruby code:
      """
      module Generic
      end
      Generic
      """
    Then there should return object
    And there should return module "Generic"
