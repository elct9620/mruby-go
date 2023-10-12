Feature: Call Function in Go
  Scenario: Can call a function defined in Go
      When I execute ruby code:
      """
      puts "Hello World"
      """
    Then there should return string "Hello World"
