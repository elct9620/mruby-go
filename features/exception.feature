Feature: Exception
  @wip
  Scenario: I can raise an exception
    When I execute ruby code:
      """
      raise "error"
      """
    Then the exception message should be "error"
