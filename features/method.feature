Feature: Method
  Scenario: I can define a method
    When I execute ruby code:
      """
      def hello
        "Hello, World!"
      end
      """
    Then there should return symbol "hello"

  Scenario: I can call the custom method
    When I execute ruby code:
      """
      def hello
        "Hello, World!"
      end
      hello
      """
    Then there should return string "Hello, World!"

  Scenario: When I call a undefined method it should raise an error
    When I execute ruby code:
      """
      hello
      """
    Then the exception message should be "undefined method 'hello' for Object"
