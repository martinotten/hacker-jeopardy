var players, categories;

/*****
  Models
*****/

var Player = Backbone.Model.extend({
  positive_score: function() {
    return this.get('score') >= 0;
  }
});

var Category = Backbone.Model.extend({
});

var Answer = Backbone.Model.extend({
});

/*****
  Collections
*****/

var Players = Backbone.Collection.extend({
  model: Player
});

var Categories = Backbone.Collection.extend({
  model: Category
});

var Answers = Backbone.Collection.extend({
  model: Answer
});

/*****
  Views
*****/

var OverviewView = Backbone.View.extend({
  collection: undefined,
  rootElement: undefined,
  source: $("#template-categories").html(),

  initialize: function(options) {
    this.collection = options['collection'];
    this.rootElement = options['root'];

    this.render();
  },

  render: function() {
    var template = Handlebars.compile(this.source);
    var categories = this.collection.map(function(category) {
      return {
        name: category.get('name'),
        answers: category.get('answers').map(function(answer) {
          return {
            answer: answer.get("answer"),
            value: answer.get("value")
          }
        }),
      };
    });
    console.log(categories);
    this.rootElement.html(template({categories: categories}));
  }
  
});

var AnswerView = Backbone.View.extend({
  answer: undefined,
  rootElement: undefined,
  template: Handlebars.compile($("#template-answer").html()),
  
  initialize: function(options) {
    this.answer = options['answer'];
    this.rootElement = options['root'];

    this.render();
  },

  render: function() {
    this.rootElement.html(this.template({answer: this.answer.get('answer')}));
  }
});

var PlayersView = Backbone.View.extend({
  collection: undefined,
  rootElement: undefined,
  source: $("#template-players").html(),

  initialize: function(options) {
    _.bindAll(this, "render");

    this.collection = options['collection'];
    this.rootElement = options['root'];
    this.listenTo(this.collection, "change:score", this.render);
    this.render();
  },

  render: function() {
    var template = Handlebars.compile(this.source);
    var players = this.collection.map(function(player) {
      return {
        name: player.get('name'),
        score: player.get('score'),
        active: player.get('active'),
        positive_score: player.positive_score()
      };
    });
    this.rootElement.html(template({players: players}));
  }
});

var AppRouter = Backbone.Router.extend({

  routes: {
    //"start"
    "overview": "overview",
    "answer/:category/:value": "answer",
    // "scorescreen",
  },

  overview: function() {
    console.log('overview');
    $.ajax({
      url: "categories.json",
      cache: false
    }).done(function(categories) {
      var categoryModels = _.map(categories, function(category) {
        var answers = _.map(category['answers'], function(answer) {
          return new Answer(answer);
        });
        category['answers'] = new Answers(answers);
        return new Category(category);
      });

      var mainElement = $("div#main");
      new OverviewView({collection: new Categories(categoryModels), root: mainElement});
    });
  },

  answer: function(category, value) {
    $.ajax({
      url: "categories.json",
      cache: false
    }).done(function(categories) {
      /*var categoryModels = _.map(categories, function(category) {
      #  var answers = _.map(category['answers'], function(answer) {
      #    return new Answer(answer);
      #  });
      #  category['answers'] = new Answers(answers);
      #  return new Category(category);
      #});
*/
      var mainElement = $("div#main");
      new AnswerView({answer: new Answer({answer: "Red Bull"}), root: mainElement});
    });
  }

});

$(document).ready(function() {
  new AppRouter();
  var playersElement = $("div#players");
  var socket = new WebSocket("ws://localhost:9090/ws");

  $.ajax({
    url: "players.json",
    cache: false
  }).done(function(players) {
    var playerModels = _.map(players, function(player) {
      return new Player(player);
    });
    new PlayersView({collection: new Players(playerModels), root: playersElement});
  });
  Backbone.history.start();
});