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
  Scenario: Int16 calculation
    When I execute ruby code:
      """
      <a> <operator> <b>
      """
    Then there should return integer <result>
    Examples:
      | a     | operator | b | result |
      | 32767 | -        | 1 | 32766  |
      | 32766 | +        | 1 | 32767  |
  Scenario: Int32 calculation
    When I execute ruby code:
      """
      <a> <operator> <b>
      """
    Then there should return integer <result>
    Examples:
      | a       | operator | b | result  |
      | 1000000 | -        | 1 | 999999  |
      | 1000000 | +        | 1 | 1000001 |
