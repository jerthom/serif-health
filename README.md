# Serif Health Takehome Solution - Jeremy Thomas

## Installing and Running the Solution
My solution is a simple Golang Cobra CLI. In order to build the executable, ensure you have Go installed on your machine, 
then simply run `go build` from the root directory. This will create an executable named `sf-takehome`. Simply running the 
executable with no arguments will default the CLI to look for a file named `2024-04-01_anthem_index.json.gz` at the root directory,
and will write the resulting output to `output.txt`. These can be configured by using the `-i` and `-o` flags respectively. 

## Investigation Process
I spent a bit more than an hour total time examining the data, coming up with a heuristic to use to filter for Anthem PPO plans in NY. While I'm not
100% satisfied with my selection of heuristic, I feel like it is suitable for the project, and with more time or background knowledge I believe
I could either refine this heuristic further or confirm its validity. 

Doing a spot check through the index file, I took random samples of EINs and fed them into the Anthem EIN lookup. I was able to find multiple rates 
files that the lookup referred to as `NY_PPO` in the human-readable link name.
(For example, EIN `04-3803249` lists `NY_PPO_GGHAMEDAP33_01_09.json.gz` through `NY_PPO_GGHAMEDAP33_09_09.json.gz`)
I inspected the web page source to copy the URLs for these links and compared them to the contents in the Index file, and found a pattern: 
all of these links seemed to share a single description: `In-Network Negotiated Rates Files`. Looking for URLs with that description, I found that
the all included a 1-2 letter state code as part of the URL, and so I formed a heuristic based on finding the description and filtering for URLs that
have the string `NY_`.

Another factor that caught my attention when I was first examining the EIN lookup site and the Index file were the links titled with the format of 
`2024-04_NY_39B0_in-network-rates_1_of_10.json.gz` in the EIN lookup. These matched up with links in the Index file that included `HighRise` in their description.
After doing some quick searching online, I came to the conclusion that HighRise is a separate entity from Anthem, although they are both subsidiaries of BCBS. 
I therefore decided to not include any of these paths in my output, although I am not 100% confident in my findings here. Additional time or background knowledge
could help raise my confidence levels on this call.

## Code Solution
I was able to code up my solution in roughly an hour, as the task of opening and streaming the gzip file was fairly straightforward. Performance-wise, my code
takes roughly 3.5 minutes to run on my machine. In order to improve the runtime of this program, my first steps would be further investigating ways to filter out
plans so that the program doesn't inspect every FileLocation in the Index. I noticed that by the half-way mark, the program stopped finding unique FileLocations, 
so I imagine there could be creative ways to leverage EINs or Plan Names to filter. I did experiment with creating a Set for EINs that the program had seen to skip over 
non-unique EINs, but there was no runtime improvements, and running a couple grep queries against the Index (`cat <index.json> | grep <EIN> | wc -l`) led me to believe 
that the EINs are unique by line. Another route to consider would be chunking the Index file and using
goroutines to concurrently scan through multiple chunks, as I imagine much of the wall-clock time of this program is consumed by IO processes. 

If I were to implement a program like this in a job setting I would also take the time to implement unit tests at the bare minimum, but I did not have time for that
during this exercise. A good example of my unit testing style can be found in my `TopTeams` repo.