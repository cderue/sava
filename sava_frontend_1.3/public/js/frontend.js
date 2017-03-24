$(document).ready(function() {
    function hipster(index, done, error) {
      $.ajax('api/message').done(function(data) {
        done(data);
      }).fail(function() {
        error({text: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."});
      });
    }

    hipster(1, function(data) {
      console.log(data);
      $('#content1').html(data.text);
    }, function(data) {
      $('#content1').html(data.text);
    });

    hipster(2, function(data) {
      console.log(data);
      $('#content2').html(data.text);
    }, function(data) {
      $('#content2').html(data.text);
    });
});
