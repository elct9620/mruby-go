Feature: Exception
  Scenario: I can raise an exception
    When I execute ruby code:
      """
      raise "error"
      """
    Then there should return exception
