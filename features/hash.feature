Feature: Hash

  Scenario: Create a new hash
    When I execute ruby code:
      """
      {"a" => 1, "b" => "str"}
      """
    Then there should return a hash
    """
    map[a:1 b:str]
    """
