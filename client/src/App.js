import React from "react";
import { BrowserRouter as Router, Switch, Route } from "react-router-dom";
import HomePage from "./pages/HomePage";
import CalendarPage from "./pages/CalendarPage";
import NewsPage from "./pages/NewsPage";
import QuotesPage from "./pages/QuotesPage";

function App() {
  return (
    <Router>
      <div>
        <Switch>
          <Route exact path="/" component={HomePage} />
          <Route path="/charts" component={QuotesPage} />
          <Route path="/calendar" component={CalendarPage} />
          <Route path="/news" component={NewsPage} />
        </Switch>
      </div>
    </Router>
  );
}

export default App;
