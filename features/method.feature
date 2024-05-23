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

  Scenario: I can call the method with argument
    When I execute ruby code:
      """
      def hello(name)
        "Hello, #{name}!"
      end
      hello("Ruby")
      """
    Then there should return string "Hello, Ruby!"

  Scenario: The missing argument raise error
    When I execute ruby code:
      """
      def hello(name)
        "Hello, #{name}!"
      end
      hello
      """
    Then the exception message should be "wrong number of arguments (given 0, expected 1)"

  Scenario: I can call the method with default argument
    When I execute ruby code:
      """
      def hello(name="World")
        "Hello, #{name}!"
      end
      hello
      """
    Then there should return string "Hello, World!"

  Scenario: I can call the method with multiple arguments
    When I execute ruby code:
      """
      def hello(name, language)
        "Hello, #{name}! I love #{language}"
      end
      hello("Ruby", "Ruby")
      """
    Then there should return string "Hello, Ruby! I love Ruby"

  Scenario: I can call the method with rest arguments
    When I execute ruby code:
      """
      def hello(*names)
        "Hello, #{names.join(', ')}!"
      end
      hello("Ruby", "Python", "JavaScript")
      """
    Then there should return string "Hello, Ruby, Python, JavaScript!"

  Scenario: When I call a undefined method it should raise an error
    When I execute ruby code:
      """
      hello
      """
    Then the exception message should be "undefined method 'hello' for Object"
