Feature: Array

  Scenario: Create a new array
    When I execute ruby code:
      """
      [1, "str", 3]
      """
    Then there should return an array
    """
    [1 str 3]
    """
