import pandas as pd 

class CDLIData:
    def __init__(self, input_path="../Data_Pipeline_go/CLDI_Extractor/new_pipeline.tsv"):
        self.input_path = input_path

    def load_data(self):
        self.data = pd.read_csv(self.input_path, sep="\t")  