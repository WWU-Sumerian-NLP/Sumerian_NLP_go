from data_loader import CDLIData
from dash import Dash, html, dcc
import plotly.express as px 

def main():
    test = CDLIData()
    test.load_data()

    app = Dash(__name__)
    fig = px.bar(test.data[567:570], x="CLDI", y="raw_translit", color="entities", barmode="group")

    app.layout = html.Div(id='parent', children=[
        html.H1(id = 'H1', children = 'Styling using html components', style = {'textAlign':'center',\
                                            'marginTop':40,'marginBottom':40}),

        dcc.Graph(
            id='example-graph',
            figure=fig
        )
    ])
    
    app.run_server(debug=True)

main()