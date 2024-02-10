Feature: Compare
  Scenario: When I compare <a> and <b> with <operator> then it should return the <res>
    When I execute ruby code:
      """
      <a> <operator> <b>
      """
    Then there should return <res>
    Examples:
      | a    | operator | b    | res   |
      | 1    | >=       | 1    | true  |
      | 1    | >        | 1    | false |
      | 1    | <        | 1    | false |
      | 1    | <=       | 1    | true  |
      | 1    | ==       | 1    | true  |
      | 2.5  | >        | 1    | true  |
      | 3.5  | <        | 1    | false |
      | 1.2  | ==       | 1    | false |
      | 3.12 | >=       | 3.11 | true  |
