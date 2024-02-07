Feature: Method
  Scenario: I can define a method
    When I execute ruby code:
      """
      def hello
        "Hello, World!"
      end
      """
    Then there should return symbol "hello"
