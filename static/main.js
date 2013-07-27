function renderPlayers(players) {
  var addPositiveScores(players) {
    return players.map(function(player) {
      player['positive_score'] = player['score'] >= 0;
      return player;
    });
  }
  players = addPositiveScores(players);
  var playersElement = $("div#players");
  var template = Handlebars.compile($("#template-players").html());
  $("div#main").html(template({players: players}));
}

function renderCategories(categories) {
  var mainElement = $("div#main");
  var template = Handlebars.compile($("#template-categories").html());
  mainElement.html(this.template({categories: categories}));
}

function renderAnswer(answer) {
  var mainElement = $("div#main");
  var template = Handlebars.compile($("#template-answer").html());
  mainElement.html(this.template({answer: answer}));
}

$(document).ready(function() {
  var socket = new WebSocket("ws://localhost:9090/ws/");
  socket.onmessage = function (event) {
    var data = $.parseJSON(event);
    if(data["players"]) {
      renderPlayers(data["players"]);
    }

    if(data["categories"]) {
      renderCategories(data["categories"]);
    }

    if(data["answer"]) {
      renderAnswer(data["answer"]);
    }
  }
});