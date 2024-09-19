Feature: Kernel
  Scenario: I can inspect Kernel module to get its name
    When I execute ruby code:
      """
      Kernel.inspect
      """
    Then there should return string "Kernel"
