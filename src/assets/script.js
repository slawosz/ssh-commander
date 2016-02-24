$("button[type=submit]").on('click', function(e) {
  e.preventDefault();

  var host = $("input[name=host]").val();
  var commands = $("textarea[name=commands]").val();
  var params = [
    {
      "Host": host,
      "Port": "10028",
      "User": user,
      "Password": password,
      "Commands": commands,
      "Prompt":">",
      "Exit":"exit"
    }
  ];

  console.log(params);
  var request = $.ajax({
    url: "http://127.0.0.1:3002/run",
    method: "POST",
    dataType: "json",
    processData: false,
    contentType: "application/json",
    data: JSON.stringify(params)
  });

  request.done(function(msg) {
    console.log(msg);
    for (var i = 0; i < msg.length; i++) {
      $("#result pre").append(msg[i].Payload.replace(/(?:\r\n|\r|\n)/g, "<br />"));
    }
    $("#myModal").modal('show');
  });
});
