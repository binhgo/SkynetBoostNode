build:
	docker login --username huynhbinh -p Cicevn2007
	docker build -t skynetboostnode .

run:
	docker run --link=SkynetMongoDB:mongodb --name SkynetBoostNode -p 6868:6868 skynetboostnode