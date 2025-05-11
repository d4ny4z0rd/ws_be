package db

import (
	"context"
	"database/sql"
	"log"
	"ws_practice_1/internal/store"
)

var dsaQuestions = []store.DSAQuestion{
	{
		Title: "Different String",
		Description: `You are given a string s consisting of lowercase English letters.

		Rearrange the characters of s to form a new string r that is not equal to s, or report that it's impossible.`,
		InputFormat: `The first line contains a single integer t (1 ≤ t ≤ 1000) — the number of test cases.

		Each of the next t lines contains a string s of length at most 10, consisting of lowercase English letters.`,
		OutputFormat: `For each test case, if it is impossible to rearrange the characters to form a different string, output "NO" (without quotes).

		Otherwise, output "YES" (without quotes) followed by the rearranged string r on the next line.

		The string r must consist of the same letters as s but must not be exactly the same as s.

		The words "YES" and "NO" are case insensitive. You may output them in any combination of uppercase and lowercase letters.

		If multiple valid answers exist, you may print any of them.`,
		ExampleInput: `
		8
		codeforces
		aaaaa
		xxxxy
		co
		d
		nutdealer
		mwistht
		hhhhhhhhhh`,
		ExampleOutput: `
		YES
		forcodesec
		NO
		YES
		xxyxx
		YES
		oc
		NO
		YES
		undertale
		YES
		thtsiwm
		NO`,
	},
	{
		Title: "Do Not Be Distracted!",
		Description: `Polycarp has 26
 		tasks. Each task is designated by a capital letter of the Latin alphabet.

		The teacher asked Polycarp to solve tasks in the following way: if Polycarp began to solve some task, then he must solve it to the end, without being distracted by another task. After switching to another task, Polycarp cannot return to the previous task.

		Polycarp can only solve one task during the day. Every day he wrote down what task he solved. Now the teacher wants to know if Polycarp followed his advice.

		For example, if Polycarp solved tasks in the following order: "DDBBCCCBBEZ", then the teacher will see that on the third day Polycarp began to solve the task 'B', then on the fifth day he got distracted and began to solve the task 'C', on the eighth day Polycarp returned to the task 'B'. Other examples of when the teacher is suspicious: "BAB", "AABBCCDDEEBZZ" and "AAAAZAAAAA".

		If Polycarp solved the tasks as follows: "FFGZZZY", then the teacher cannot have any suspicions. Please note that Polycarp is not obligated to solve all tasks. Other examples of when the teacher doesn't have any suspicious: "BA", "AFFFCC" and "YYYYY".

		Help Polycarp find out if his teacher might be suspicious.`,
		InputFormat: `The first line contains an integer t
 		(1≤t≤1000). Then t test cases follow.

		The first line of each test case contains one integer n
 		(1≤n≤50) — the number of days during which Polycarp solved tasks.

		The second line contains a string of length n, consisting of uppercase Latin letters, which is the order in which Polycarp solved the tasks.`,
		OutputFormat: `For each test case output:

		"YES", if the teacher cannot be suspicious;
		"NO", otherwise.
		You may print every letter in any case you want (so, for example, the strings yEs, yes, Yes and YES are all recognized as positive answer).`,
		ExampleInput: `
		5
		3
		ABA
		11
		DDBBCCCBBEZ
		7
		FFGZZZY
		1
		Z
		2
		AB
		`,
		ExampleOutput: `
		NO
		NO
		YES
		YES
		YES
		`,
	},
	{
		Title: "Good Kid",
		Description: `Slavic is preparing a present for a friend's birthday. He has an array a of n
		digits and the present will be the product of all these digits. Because Slavic is a good kid who wants to make the biggest product possible, he wants to add 1
		to exactly one of his digits. What is the maximum product Slavic can make?`,
		InputFormat: `The first line contains a single integer t
 		(1≤t≤104) — the number of test cases.

		The first line of each test case contains a single integer n
		(1≤n≤9
		) — the number of digits.

		The second line of each test case contains n
		space-separated integers ai
		(0≤ai≤9
		) — the digits in the array.

		`,
		OutputFormat: `For each test case, output a single integer — the maximum product Slavic can make, by adding 1
		to exactly one of his digits.`,
		ExampleInput: `
		4
		4
		2 2 1 2
		3
		0 1 2
		5
		4 3 2 3 4
		9
		9 9 9 9 9 9 9 9 9
		`,
		ExampleOutput: `
		16
		2
		432
		430467210
		`,
	},
	{
		Title: "Fair Division",
		Description: `Alice and Bob received n
 		candies from their parents. Each candy weighs either 1 gram or 2 grams. Now they want to divide all candies among themselves fairly so that the total weight of Alice's candies is equal to the total weight of Bob's candies.

		Check if they can do that.

		Note that candies are not allowed to be cut in half.`,
		InputFormat: `The first line contains one integer t
		(1≤t≤104
		) — the number of test cases. Then t
		test cases follow.

		The first line of each test case contains an integer n
		(1≤n≤100
		) — the number of candies that Alice and Bob received.

		The next line contains n
		integers a1,a2,…,an
		— the weights of the candies. The weight of each candy is either 1
		or 2
		.

		It is guaranteed that the sum of n
		over all test cases does not exceed 105
		.
		`,
		OutputFormat: `For each test case, output on a separate line:

		"YES", if all candies can be divided into two sets with the same weight;
		"NO" otherwise.
		You can output "YES" and "NO" in any case (for example, the strings yEs, yes, Yes and YES will be recognized as positive).`,
		ExampleInput: `
		5
		2
		1 1
		2
		1 2
		4
		1 2 1 2
		3
		2 2 2
		3
		2 1 2

		`,
		ExampleOutput: `
		YES
		NO
		YES
		NO
		NO
		`,
	},
	{
		Title: "Polycarp and Coins",
		Description: `Polycarp must pay exactly n
		burles at the checkout. He has coins of two nominal values: 1
		burle and 2
		burles. Polycarp likes both kinds of coins equally. So he doesn't want to pay with more coins of one type than with the other.

		Thus, Polycarp wants to minimize the difference between the count of coins of 1
		burle and 2
		burles being used. Help him by determining two non-negative integer values c1
		and c2
		which are the number of coins of 1
		burle and 2
		burles, respectively, so that the total value of that number of coins is exactly n
		(i. e. c1+2⋅c2=n
		), and the absolute value of the difference between c1
		and c2
		is as little as possible (i. e. you must minimize |c1−c2|
		).

		Note that candies are not allowed to be cut in half.`,
		InputFormat: `The first line contains one integer t
		(1≤t≤104
		) — the number of test cases. Then t
		test cases follow.

		Each test case consists of one line. This line contains one integer n
		(1≤n≤109
		) — the number of burles to be paid by Polycarp.
		`,
		OutputFormat: `For each test case, output a separate line containing two integers c1
		and c2
		(c1,c2≥0
		) separated by a space where c1
		is the number of coins of 1
		burle and c2
		is the number of coins of 2
		burles. If there are multiple optimal solutions, print any one.`,
		ExampleInput: `
		6
		1000
		30
		1
		32
		1000000000
		5
		`,
		ExampleOutput: `
		334 333
		10 10
		1 0
		10 11
		333333334 333333333
		1 2
		`,
	},
}

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()

	tx, _ := db.BeginTx(ctx, nil)

	for _, q := range dsaQuestions {
		if err := store.Questions.Create(ctx, &q); err != nil {
			_ = tx.Rollback()
			log.Println("Error creating user:", err)
			return
		}
	}

	tx.Commit()

	log.Println("Seeding complete")
}
