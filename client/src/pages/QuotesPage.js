import React from "react";
import NavigationMenu from "../components/NavigationMenu/NavigationMenu";
import QuoteSelector from "../components/QuoteSelector/QuoteSelector";

function QuotesPage() {
  return (
    <React.Fragment>
      <NavigationMenu />
      <QuoteSelector />
    </React.Fragment>
  );
}

export default QuotesPage;
