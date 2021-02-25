import React from "react";
import { BrowserRouter as Router, Switch, Route } from "react-router-dom";
import HomePage from "./pages/HomePage";
import ChartsPage from "./pages/ChartsPage";
import CalendarPage from "./pages/CalendarPage";
import NewsPage from "./pages/NewsPage";
import QuotesPage from "./pages/QuotesPage";

function App() {
  return (
    <Router>
      <div>
        <Switch>
          <Route exact path="/" component={HomePage} />
          <Route path="/charts" component={ChartsPage} />
          <Route path="/calendar" component={CalendarPage} />
          <Route path="/news" component={NewsPage} />
          <Route path="/quotes" component={QuotesPage} />
        </Switch>
      </div>
    </Router>
  );
}

export default App;
