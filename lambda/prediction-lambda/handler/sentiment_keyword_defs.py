SENTIMENT_KEYWORDS = {
        "USD": {
            "positive": [
                "usd/",
                "u.s.",
                "greenback",
                "buck",
                "barnie",
                "america",
                "united states",
            ],
            "negative": ["/usd", "cable"],
        },
        "AUD": {
            "positive": ["aud/", "gold", "aussie", "australia"],
            "negative": ["/aud"],
        },
        "GBP": {
            "positive": [
                "gbp/",
                "sterling",
                "pound",
                "u.k.",
                "united kingdom",
                "cable",
                "guppy",
            ],
            "negative": ["/gbp"],
        },
        "NZD": {
            "positive": ["nzd/", "gold", "kiwi", "new zealand"],
            "negative": ["/nzd"],
        },
        "CAD": {"positive": ["cad/", "oil", "loonie", "canada"], "negative": ["/cad"]},
        "CHF": {"positive": ["chf/", "swiss"], "negative": ["/chf"]},
        "JPY": {"positive": ["jpy/", "asian", "japan"], "negative": ["/jpy", "guppy"]},
        "EUR": {"positive": ["eur/", "fiber", "euro"], "negative": ["/eur"]},
    }