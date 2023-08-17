Feature: Math
  Scenario: Small number calculation
    When I execute ruby code:
      """
      <a> <operator> <b>
      """
    Then there should return integer <result>
    Examples:
      | a | operator | b | result |
      | 1 | -        | 2 | -1     |
      | 1 | -        | 1 | 0      |
      | 2 | -        | 1 | 1      |
      | 1 | +        | 1 | 2      |
      | 2 | +        | 1 | 3      |
      | 3 | +        | 1 | 4      |
      | 2 | +        | 3 | 5      |
      | 4 | +        | 2 | 6      |
      | 3 | +        | 4 | 7      |
