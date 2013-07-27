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
  
});

var AnswerView = Backbone.View.extend({
  
});

var PlayersView = Backbone.View.extend({
  
});

var AppRouter = Backbone.Router.extend({

  routes: {
    //"start"
    "overview": "overview",
    "answer/:category/:value": "answer",
    // "scorescreen",
  },

  overview: function() {
    var mainElement = $("div#main");
    new OverviewView(mainElement, players, categories);
  },

  answer: function(category, value) {
    new AnswerView();
  }

});

$(document).ready(function() {
  new AppRouter();
  var playersElement = $("div#players");
  new PlayersView();
  Backbone.history.start({pushState: true});
});