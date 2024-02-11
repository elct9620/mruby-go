Feature: Condition
  Scenario: When condition is <condition> then it should return <expected>
    When I execute ruby code:
      """
      if <condition>
        "yes"
      else
        "no"
      end
      """
    Then there should return string "<expected>"
    Examples:
      | condition | expected |
      | 1 > 2     | no       |
      | 1 < 2     | yes      |
