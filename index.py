#require python3
#pip3 install tabula
#pip3 install PdfMiner3K

import tabula
from pdfminer.pdfparser import PDFParser, PDFDocument
from pdfminer.pdfinterp import PDFResourceManager, PDFPageInterpreter
from pdfminer.converter import PDFPageAggregator
from pdfminer.layout import LAParams, LTTextBox, LTTextLine
import os
import platform
import sys

def parse_pdf(path, output_path):
    with open(path, 'rb') as fp:
        parser = PDFParser(fp)
        doc = PDFDocument()
        parser.set_document(doc)
        doc.set_parser(parser)
        doc.initialize('')
        rsrcmgr = PDFResourceManager()
        laparams = LAParams()
        laparams.char_margin = 1.0
        laparams.word_margin = 1.0
        device = PDFPageAggregator(rsrcmgr, laparams=laparams)
        interpreter = PDFPageInterpreter(rsrcmgr, device)
        extracted_text = ''
        for page in doc.get_pages():
            interpreter.process_page(page)
            layout = device.get_result()
            for lt_obj in layout:
                if isinstance(lt_obj, LTTextBox) or isinstance(lt_obj, LTTextLine):
                    extracted_text += lt_obj.get_text()
    with open(output_path, "w", encoding="utf-8") as f:
        f.write(extracted_text)

def list_all_files(rootdir):
    _files = []
    list = os.listdir(rootdir) 
    for i in range(0,len(list)):
           path = os.path.join(rootdir,list[i])
           if os.path.isdir(path):
              _files.extend(list_all_files(path))
           if os.path.isfile(path):
              _files.append(path)
    return _files

def get_file_name(filename):
    rawname=filename.split(".")
    if len(rawname)<=1:
        return "unkonw_pdf"
    return rawname[0]

if __name__ == "__main__":
    title=sys.argv[1]
    filename=sys.argv[2]
    srcDir=""
    tarDir=""
    platStr=platform.system()
    realPath=sys.path[0]
    if "Windows"==platStr:
        #srcDir="E:\\royalchen\\gopath\\src\\wx\\xcx_yaohao_server\\data\\"+title+"\\pdf\\"
        #tarDir="E:\\royalchen\\gopath\\src\\wx\\xcx_yaohao_server\\data\\"+title+"\\txt\\"
        srcDir=realPath+"\\data\\"+title+"\\pdf\\"
        tarDir=realPath+"\\data\\"+title+"\\txt\\"
    else:
        srcDir="./data/"+title+"/pdf/"
        tarDir="./data/"+title+"/txt/"
    print("start transform PDF to txt file in ",srcDir,",filename:",filename)
    name=get_file_name(filename)
    newfilename=tarDir+name+".txt"
    tarpdf=srcDir+filename
    print("start parse file:",tarpdf," ==========>",newfilename)
    parse_pdf(tarpdf,newfilename)
    print("parse file ",name," complete")
   