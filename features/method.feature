Feature: Method
  @wip
  Scenario: I can define a method
    When I execute ruby code:
      """
      def hello
        "Hello, World!"
      end
      puts hello
      """
    Then there should return string "Hello World"
