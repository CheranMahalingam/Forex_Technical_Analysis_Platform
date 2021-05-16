import React from "react";
import { makeStyles } from "@material-ui/core/styles";
import Typography from "@material-ui/core/Typography";

const useStyles = makeStyles(() => ({
    root: {
        color: "#FFFFFF",
        marginTop: 100,
        marginLeft: "15%"
    }
}))

function EmptyQuote() {
    const classes = useStyles();

    return (
        <React.Fragment>
            <Typography variant="body1" className={classes.root}>
                No currency pair selected :(
            </Typography>
        </React.Fragment>
    )
}

export default EmptyQuote;
