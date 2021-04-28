import sys
import re
import os
import time
import subprocess
from glob import glob
from multiprocessing.dummy import Pool as ThreadPool 

GO_SRC_PATH = os.environ['GOPATH'] + '/src/'
PATH = GO_SRC_PATH + "github.com/guilhermeCoutinho/api-studies"
goFiles = [y for x in os.walk(PATH) for y in glob(os.path.join(x[0], '*.go')) if "/dal/" in y or "/services/" in y or "/server/http/" in y ]
pattern = re.compile("^type [A-Z][^ .]* interface {")
currentProgress = 0
targetProgress = 0

def main():
	global targetProgress
	fileInterfaceTuples = SearchAllFiles()
	targetProgress = len (fileInterfaceTuples)
	displayProgress()
	pool = ThreadPool(8)
	out = pool.map(proccess, fileInterfaceTuples)
	print ("")
	for s in out:
		if s != "":
			print (s)
	pool.close()

def proccess (fileInterfaceTuple):
	global currentProgress
	filePath = fileInterfaceTuple[0]
	interfaceName = fileInterfaceTuple[1]
	packagePath = extractPackagePath(filePath)
	out = ""
	if shouldGenerateMock(filePath,interfaceName):
		runMockGen(packagePath,interfaceName,pascal_case_to_snake_case(interfaceName))
		out = interfaceName + " exported to " + "mocks/" + pascal_case_to_snake_case(interfaceName)
	currentProgress += 1
	displayProgress()
	return out

def shouldGenerateMock (fileName, interface):
	mockFilePath = PATH + '/mocks/' + pascal_case_to_snake_case(interface) + '.go'
	mockExists = os.path.isfile(mockFilePath)
	if not mockExists:
		return True
	mockModificationTime = os.path.getmtime(mockFilePath)
	fileToBeMockedModificationTime = os.path.getmtime(fileName)
	return mockModificationTime < fileToBeMockedModificationTime

def runMockGen (packagePath , interfaceName , outputFileName) :
	packageName = packagePath.split('/')[-1]
	command = "mockgen " + packagePath + " " + interfaceName + "| sed 's/mock_" + packageName + "/mocks/' > mocks/" + outputFileName + ".go"
	subprocess.call(command,shell=True)

def displayProgress ():
	global currentProgress
	global targetProgress
	progress =  "\r[" + ('=' * currentProgress ) + ('.' * (targetProgress-currentProgress) ) + "]"
	sys.stdout.write(progress)
	sys.stdout.flush()

def extractPackagePath (filePath):
	filePath = filePath.replace(GO_SRC_PATH,'')
	filePath = '/'.join( filePath.split('/')[:-1] )
	return filePath

def SearchAllFiles ():
	fileInterfaceTuples = []
	for i in goFiles:
		for fileInterfaceTuple in SeachForInterfacesInFile(i):
			fileInterfaceTuples += [fileInterfaceTuple] if fileInterfaceTuple is not None else []
	return fileInterfaceTuples

def SeachForInterfacesInFile (fileName):
	interfaces = []
	for i, line in enumerate(open(fileName)):
	    for match in re.finditer(pattern, line):
	        interfaces.append( (fileName,line.split()[1]) )
	return interfaces

def pascal_case_to_snake_case (name):
    s1 = re.sub('(.)([A-Z][a-z]+)', r'\1_\2', name)
    return re.sub('([a-z0-9])([A-Z])', r'\1_\2', s1).lower()

if __name__ == "__main__":
	main()