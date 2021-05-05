import React, {useState} from "react";
import Backdrop from "@material-ui/core/Backdrop";
import CircularProgress from "@material-ui/core/CircularProgress";
import { makeStyles } from "@material-ui/core/styles";
import NavigationMenu from "../components/NavigationMenu/NavigationMenu";
import QuoteSelector from "../components/QuoteSelector/QuoteSelector";

const useStyles = makeStyles((theme) => ({
  backdrop: {
    zIndex: theme.zIndex.drawer + 1,
    color: "#fff",
  },
}));

function QuotesPage() {
  const [loading, setLoading] = useState(true);

  const classes = useStyles();

  const handleWebsocketOpen = () => {
    setLoading(false);
  }

  return (
    <React.Fragment>
      <Backdrop className={classes.backdrop} open={loading}>
        <CircularProgress color="inherit" />
      </Backdrop>
      <NavigationMenu />
      <QuoteSelector handleEndLoading={handleWebsocketOpen}/>
    </React.Fragment>
  );
}

export default QuotesPage;
