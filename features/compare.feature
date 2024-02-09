Feature: Compare
  Scenario: When I compare two numbers then it should return the result
    When I execute ruby code:
      """
      <a> <operator> <b>
      """
    Then there should return <res>
    Examples:
      | a | operator | b | res   |
      | 1 | >=       | 1 | true  |
      | 1 | >        | 1 | false |
      | 1 | <        | 1 | false |
      | 1 | <=       | 1 | true  |
      | 1 | ==       | 1 | true  |
