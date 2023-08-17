Feature: Boolean
  Scenario: Can return boolean
    When I execute ruby code:
      """
      <bool>
      """
    Then there should return <bool>
    Examples:
      | bool  |
      | true  |
      | false |
