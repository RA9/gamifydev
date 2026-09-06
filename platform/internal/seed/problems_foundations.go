package seed

import "github.com/RA9/gamifydev/platform/internal/store"

// foundationProblems are CLI-style C exercises embedded directly into the
// Foundations courses. Their stdin/stdout contracts let the same sandbox that
// grades C checkpoints judge them without exposing a language-specific harness.
var foundationProblems = []seedProblem{
	{
		slug: "signal-checksum", title: "Signal checksum", difficulty: "easy", topic: "C Fundamentals",
		statement: `A rover sends a packet as a count followed by that many signed integers.
Read the packet and print the sum of its values modulo 100. Keep the result
non-negative, even when the packet contains negative values.

Input:

    4
    30 40 -5 50

Output:

    15`,
		timeLimitMs: 3000, workloadMinutes: 100, mode: store.ModeFun, language: "c", solutionLang: "c",
		starters: cStarter(`#include <stdio.h>

int main(void) {
    int n;
    if (scanf("%d", &n) != 1) return 1;
    /* Read n values and print the non-negative checksum. */
    return 0;
}
`),
		tests: []store.Check{
			{Label: "checksums a mixed packet", Stdin: "4\n30 40 -5 50\n", Test: `_out.strip() == "15"`, Points: 1},
			{Label: "handles an empty packet", Stdin: "0\n", Test: `_out.strip() == "0"`, Points: 1},
			{Label: "keeps a negative total non-negative", Stdin: "2\n-60 -55\n", Test: `_out.strip() == "85"`, Hidden: true, Points: 1},
			{Label: "reduces values before they can overflow the sum", Stdin: "2\n1 2147483647\n", Test: `_out.strip() == "48"`, Hidden: true, Points: 1},
		},
		solution: `#include <stdio.h>
int main(void){int n,x,sum=0;if(scanf("%d",&n)!=1)return 1;for(int i=0;i<n;i++){if(scanf("%d",&x)!=1)return 1;sum=(sum+(x%100))%100;}if(sum<0)sum+=100;printf("%d\n",sum);return 0;}
`,
	},
	{
		slug: "stack-command-trace", title: "Stack command trace", difficulty: "easy", topic: "Data Structures",
		statement: `Implement a fixed-capacity integer stack. The first input line is the
number of commands. Each following command is either ~PUSH value~ or ~POP~.
Print the value removed by every pop, or ~EMPTY~ when the stack has no value.
The tests never exceed 1,000 stored values.`,
		timeLimitMs: 3000, workloadMinutes: 70, mode: store.ModePractical, language: "c", solutionLang: "c",
		starters: cStarter(`#include <stdio.h>
#include <string.h>

int main(void) {
    int stack[1000], top = 0, commands;
    if (scanf("%d", &commands) != 1) return 1;
    /* Process PUSH and POP commands. */
    return 0;
}
`),
		tests: []store.Check{
			{Label: "uses last-in, first-out order", Stdin: "5\nPUSH 3\nPUSH 8\nPOP\nPUSH 2\nPOP\n", Test: `_out.strip().splitlines() == ["8", "2"]`, Points: 1},
			{Label: "reports an empty stack", Stdin: "2\nPOP\nPOP\n", Test: `_out.strip().splitlines() == ["EMPTY", "EMPTY"]`, Points: 1},
			{Label: "retains values below the top", Stdin: "6\nPUSH -4\nPUSH 7\nPUSH 9\nPOP\nPOP\nPOP\n", Test: `_out.strip().splitlines() == ["9", "7", "-4"]`, Hidden: true, Points: 1},
		},
		solution: `#include <stdio.h>
#include <string.h>
int main(void){int stack[1000],top=0,n,value;char command[8];if(scanf("%d",&n)!=1)return 1;for(int i=0;i<n;i++){if(scanf("%7s",command)!=1)return 1;if(strcmp(command,"PUSH")==0){if(scanf("%d",&value)!=1||top>=1000)return 1;stack[top++]=value;}else if(strcmp(command,"POP")==0){if(top==0)puts("EMPTY");else printf("%d\n",stack[--top]);}}return 0;}
`,
	},
	{
		slug: "pair-sum-stream", title: "Pair sum stream", difficulty: "medium", topic: "Complexity",
		statement: `Read ~n~, a target, and ~n~ integers. Print ~YES~ when two different
positions add to the target and ~NO~ otherwise. Design for thousands of values;
a quadratic comparison of every pair is not the intended solution.`,
		timeLimitMs: 3000, workloadMinutes: 70, mode: store.ModePractical, language: "c", solutionLang: "c",
		starters: cStarter(`#include <stdio.h>
#include <stdlib.h>

int main(void) {
    int n, target;
    if (scanf("%d %d", &n, &target) != 2) return 1;
    /* Read, organize, and search the values. */
    return 0;
}
`),
		tests: []store.Check{
			{Label: "finds an existing pair", Stdin: "4 9\n2 7 11 15\n", Test: `_out.strip() == "YES"`, Points: 1},
			{Label: "does not reuse one position", Stdin: "1 10\n5\n", Test: `_out.strip() == "NO"`, Points: 1},
			{Label: "handles negative values", Stdin: "5 -2\n8 -3 1 7 20\n", Test: `_out.strip() == "YES"`, Hidden: true, Points: 1},
		},
		solution: `#include <stdio.h>
#include <stdlib.h>
static int cmp(const void*a,const void*b){int x=*(const int*)a,y=*(const int*)b;return (x>y)-(x<y);}int main(void){int n,target;if(scanf("%d %d",&n,&target)!=2||n<0)return 1;int*xs=malloc((size_t)n*sizeof(int));if(n>0&&!xs)return 1;for(int i=0;i<n;i++)if(scanf("%d",&xs[i])!=1){free(xs);return 1;}qsort(xs,(size_t)n,sizeof(int),cmp);int l=0,r=n-1,found=0;while(l<r){long sum=(long)xs[l]+xs[r];if(sum==target){found=1;break;}if(sum<target)l++;else r--;}puts(found?"YES":"NO");free(xs);return 0;}
`,
	},
	{
		slug: "route-hop-count", title: "Route hop count", difficulty: "medium", topic: "Graphs",
		statement: `A directed network has nodes numbered from 0. Read ~n m~, then ~m~
edges, then a start and goal. Print the fewest number of edges needed to reach
the goal, or ~-1~ when it is unreachable. Cycles are allowed.`,
		timeLimitMs: 3000, workloadMinutes: 70, mode: store.ModePractical, language: "c", solutionLang: "c",
		starters: cStarter(`#include <stdio.h>
#include <stdlib.h>

int main(void) {
    int n, m;
    if (scanf("%d %d", &n, &m) != 2) return 1;
    /* Build the graph and use breadth-first search. */
    return 0;
}
`),
		tests: []store.Check{
			{Label: "takes the direct route", Stdin: "4 4\n0 1\n1 2\n2 3\n0 3\n0 3\n", Test: `_out.strip() == "1"`, Points: 1},
			{Label: "reports an unreachable goal", Stdin: "3 1\n0 1\n0 2\n", Test: `_out.strip() == "-1"`, Points: 1},
			{Label: "handles a cycle", Stdin: "4 4\n0 1\n1 0\n1 2\n2 3\n0 3\n", Test: `_out.strip() == "3"`, Hidden: true, Points: 1},
		},
		solution: `#include <stdio.h>
#include <stdlib.h>
int main(void){int n,m;if(scanf("%d %d",&n,&m)!=2||n<=0||m<0)return 1;unsigned char*adj=calloc((size_t)n*n,1);int*dist=malloc((size_t)n*sizeof(int));int*q=malloc((size_t)n*sizeof(int));if(!adj||!dist||!q){free(adj);free(dist);free(q);return 1;}for(int i=0,u,v;i<m;i++){if(scanf("%d %d",&u,&v)!=2||u<0||u>=n||v<0||v>=n)return 1;adj[u*n+v]=1;}int start,goal;if(scanf("%d %d",&start,&goal)!=2)return 1;for(int i=0;i<n;i++)dist[i]=-1;int head=0,tail=0;q[tail++]=start;dist[start]=0;while(head<tail){int u=q[head++];for(int v=0;v<n;v++)if(adj[u*n+v]&&dist[v]<0){dist[v]=dist[u]+1;q[tail++]=v;}}printf("%d\n",dist[goal]);free(adj);free(dist);free(q);return 0;}
`,
	},
	{
		slug: "count-set-bits", title: "Count set bits", difficulty: "easy", topic: "Computer Architecture",
		statement: `Read one unsigned 32-bit integer and print how many bits are 1 in its
binary representation. Use bit operations rather than converting the number to
text.`,
		timeLimitMs: 3000, workloadMinutes: 70, mode: store.ModePractical, language: "c", solutionLang: "c",
		starters: cStarter(`#include <stdint.h>
#include <inttypes.h>
#include <stdio.h>

int main(void) {
    uint32_t value;
    if (scanf("%" SCNu32, &value) != 1) return 1;
    /* Count the set bits and print the count. */
    return 0;
}
`),
		tests: []store.Check{
			{Label: "counts alternating bits", Stdin: "10\n", Test: `_out.strip() == "2"`, Points: 1},
			{Label: "zero has no set bits", Stdin: "0\n", Test: `_out.strip() == "0"`, Points: 1},
			{Label: "counts every bit in a full word", Stdin: "4294967295\n", Test: `_out.strip() == "32"`, Hidden: true, Points: 1},
		},
		solution: `#include <stdint.h>
#include <inttypes.h>
#include <stdio.h>
int main(void){uint32_t value;if(scanf("%" SCNu32,&value)!=1)return 1;int count=0;while(value){value&=value-1;count++;}printf("%d\n",count);return 0;}
`,
	},
	{
		slug: "permission-mask", title: "Permission mask", difficulty: "easy", topic: "Linux",
		statement: `Unix permission digits encode read as 4, write as 2, and execute as
1. Read one digit from 0 through 7 and print exactly three characters: ~r~ or
~-~, then ~w~ or ~-~, then ~x~ or ~-~.`,
		timeLimitMs: 3000, workloadMinutes: 70, mode: store.ModePractical, language: "c", solutionLang: "c",
		starters: cStarter(`#include <stdio.h>

int main(void) {
    int mode;
    if (scanf("%d", &mode) != 1) return 1;
    /* Decode the permission bits. */
    return 0;
}
`),
		tests: []store.Check{
			{Label: "decodes full access", Stdin: "7\n", Test: `_out.strip() == "rwx"`, Points: 1},
			{Label: "decodes read-only access", Stdin: "4\n", Test: `_out.strip() == "r--"`, Points: 1},
			{Label: "decodes write and execute", Stdin: "3\n", Test: `_out.strip() == "-wx"`, Hidden: true, Points: 1},
		},
		solution: `#include <stdio.h>
int main(void){int mode;if(scanf("%d",&mode)!=1||mode<0||mode>7)return 1;putchar(mode&4?'r':'-');putchar(mode&2?'w':'-');putchar(mode&1?'x':'-');putchar('\n');return 0;}
`,
	},
}

var foundationCourseProblems = map[string][]string{
	"c":                       {"signal-checksum"},
	"data_structures":         {"stack-command-trace"},
	"complexity_and_analysis": {"pair-sum-stream"},
	"algorithms":              {"route-hop-count"},
	"how_computers_work":      {"count-set-bits"},
	"linux":                   {"permission-mask"},
}
