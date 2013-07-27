var players, categories;

var OverviewView = Backbone.View.extend({
  
});

var QuestionView = Backbone.View.extend({
  
});

var PlayersView = Backbone.View.extend({
  
});

PlayersView.prototype.views = {};
PlayersView.prototype.views['Player'] = Backbone.View.extend({
  
});



var AppRouter = Backbone.Router.extend({

  routes: {
    //"start"
    "overview": "overview",
    "question/:category_name/:question_id":        "question",  // #search/kiwis
    // "scorescreen",
  },

  overview: function() {
    var mainElement = $("div#main");
    new OverviewView(mainElement, players, categories);
  },

  question: function(category_name, question_id) {
    
  }

});

$(document).ready(function() {
  new AppRouter();
  Backbone.history.start({pushState: true});
});