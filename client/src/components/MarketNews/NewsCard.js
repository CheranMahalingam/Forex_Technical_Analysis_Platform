import React from "react";
import { makeStyles } from "@material-ui/core/styles";
import Card from "@material-ui/core/Card";
import CardActions from "@material-ui/core/CardActions";
import CardHeader from "@material-ui/core/CardHeader";
import CardMedia from "@material-ui/core/CardMedia";
import CardContent from "@material-ui/core/CardContent";
import Typography from "@material-ui/core/Typography";
import Button from "@material-ui/core/Button";
import DefaultImage from "../../assets/forex_broken_link.jpg";

const useStyles = makeStyles(() => ({
  card: {
    marginTop: 20,
    marginLeft: 15,
    marginRight: 15,
    padding: 5,
  },
  source: {
      marginLeft: 5
  },
  header: {
    textAlign: "center",
    marginTop: -10
  },
  media: {
    textAlign: "center",
  },
  readMore: {
    float: "right",
    marginRight: 20,
    marginTop: -20,
  },
  date: {
      color: "#4C4E52",
      marginTop: -10
  }
}));

function NewsCard(props) {
  const classes = useStyles();

  return (
    <React.Fragment>
      <Card className={classes.card}>
        <Typography className={classes.source} variant="h6">
          {props.source}
        </Typography>
        <CardHeader className={classes.header} title={props.title} />
        <CardMedia className={classes.media}>
          <img
            src={props.image}
            width="90%"
            onError={(event) => {
              event.target.onerror = null;
              event.target.src = DefaultImage;
            }}
          />
        </CardMedia>
        <CardContent>
          <Typography className={classes.date} variant="body2">
            {props.timestamp}
          </Typography>
          <Typography variant="body2">{props.summary}</Typography>
        </CardContent>
        <CardActions className={classes.readMore}>
          <a href={props.newsurl}>
            <Button size="small" color="primary">
              Read More
            </Button>
          </a>
        </CardActions>
      </Card>
    </React.Fragment>
  );
}

export default NewsCard;
